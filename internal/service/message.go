package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/TranTheTuan/vna/configs"
	"github.com/TranTheTuan/vna/internal/domain"
	"github.com/TranTheTuan/vna/internal/repository"
)

// Sentinel errors returned by MessageService methods.
var (
	ErrEmptyMessage    = errors.New("service.message: message cannot be empty")
	ErrUpstreamFailed  = errors.New("service.message: upstream API returned an error")
	ErrUpstreamTimeout = errors.New("service.message: upstream API timed out")
	ErrInvalidLimit    = errors.New("service.message: limit must be between 1 and 100")
)

// MessageService defines operations for sending and listing chat messages.
type MessageService interface {
	Send(ctx context.Context, userID, question string) (*domain.Message, error)
	// SendStream calls the OpenResponses API with stream:true, calls onDelta for
	// each incremental text chunk, saves the full answer, and returns the message.
	SendStream(ctx context.Context, userID, question string, onDelta func(chunk string)) (*domain.Message, error)
	List(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error)
}

// openResponsesRequest is the body sent to POST /v1/responses.
type openResponsesRequest struct {
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
	User   string `json:"user"`
	Input  string `json:"input"`
}

// openResponsesResp is the expected response shape from the non-streaming OpenResponses API.
type openResponsesResp struct {
	Output []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

// sseDeltaData is the payload for "response.output_text.delta" SSE events.
type sseDeltaData struct {
	Delta string `json:"delta"`
}

type messageService struct {
	cfg        *configs.Config
	repo       repository.MessageRepository
	httpClient *http.Client
	logger     *slog.Logger
}

// NewMessageService creates a MessageService. The HTTP client has no global timeout;
// callers control cancellation via context (streaming requires long-lived connections).
func NewMessageService(cfg *configs.Config, repo repository.MessageRepository, logger *slog.Logger) MessageService {
	return &messageService{
		cfg:        cfg,
		repo:       repo,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

// Send calls the OpenResponses API (non-streaming) with the user's question, saves the Q&A, and returns it.
func (s *messageService) Send(ctx context.Context, userID, question string) (*domain.Message, error) {
	if question == "" {
		return nil, ErrEmptyMessage
	}

	answer, err := s.callOpenResponses(ctx, userID, question)
	if err != nil {
		s.logger.Error("OpenResponses call failed", "error", err, "userID", userID)
		return nil, err
	}

	msg, err := s.repo.Save(ctx, &domain.Message{
		UserID:   userID,
		Question: question,
		Answer:   answer,
	})
	if err != nil {
		s.logger.Error("save message failed", "error", err, "userID", userID)
		return nil, fmt.Errorf("service.message: save message: %w", err)
	}
	return msg, nil
}

// SendStream calls the OpenResponses API with stream:true, invokes onDelta for each
// text chunk received, saves the full accumulated answer, and returns the saved message.
func (s *messageService) SendStream(ctx context.Context, userID, question string, onDelta func(chunk string)) (*domain.Message, error) {
	if question == "" {
		return nil, ErrEmptyMessage
	}

	answer, err := s.streamOpenResponses(ctx, userID, question, onDelta)
	if err != nil {
		s.logger.Error("OpenResponses stream failed", "error", err, "userID", userID)
		return nil, err
	}

	msg, err := s.repo.Save(ctx, &domain.Message{
		UserID:   userID,
		Question: question,
		Answer:   answer,
	})
	if err != nil {
		s.logger.Error("save message failed", "error", err, "userID", userID)
		return nil, fmt.Errorf("service.message: save message: %w", err)
	}
	return msg, nil
}

// List returns a paginated slice of messages for the given user.
// limit is clamped to [1, 100]; defaults to 20 when 0 is passed.
func (s *messageService) List(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error) {
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 {
		return nil, "", ErrInvalidLimit
	}
	return s.repo.ListByUser(ctx, userID, limit, cursor)
}

// callOpenResponses sends a non-streaming POST /v1/responses request and extracts the answer text.
func (s *messageService) callOpenResponses(ctx context.Context, userID, question string) (string, error) {
	payload := openResponsesRequest{
		Model:  s.cfg.ChatServer.Model,
		Stream: false,
		User:   userID,
		Input:  question,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("service.message: marshal request: %w", err)
	}

	url := s.cfg.ChatServer.BaseUrl + "/v1/responses"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("service.message: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.ChatServer.AuthToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", ErrUpstreamTimeout
		}
		return "", fmt.Errorf("%w: %v", ErrUpstreamFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d", ErrUpstreamFailed, resp.StatusCode)
	}

	var parsed openResponsesResp
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("%w: failed to parse response JSON", ErrUpstreamFailed)
	}

	if len(parsed.Output) == 0 || len(parsed.Output[0].Content) == 0 {
		return "", fmt.Errorf("%w: unexpected empty output", ErrUpstreamFailed)
	}
	answer := parsed.Output[0].Content[0].Text
	if answer == "" {
		return "", fmt.Errorf("%w: empty answer text in response", ErrUpstreamFailed)
	}

	return answer, nil
}

// streamOpenResponses sends a streaming POST /v1/responses request (stream:true),
// reads SSE events line-by-line, calls onDelta for each response.output_text.delta chunk,
// and returns the fully accumulated answer string.
func (s *messageService) streamOpenResponses(ctx context.Context, userID, question string, onDelta func(chunk string)) (string, error) {
	payload := openResponsesRequest{
		Model:  s.cfg.ChatServer.Model,
		Stream: true,
		User:   userID,
		Input:  question,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("service.message: marshal request: %w", err)
	}

	url := s.cfg.ChatServer.BaseUrl + "/v1/responses"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("service.message: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.ChatServer.AuthToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", ErrUpstreamTimeout
		}
		return "", fmt.Errorf("%w: %v", ErrUpstreamFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d", ErrUpstreamFailed, resp.StatusCode)
	}

	// Use a 1MB scanner buffer to handle large SSE data lines.
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	var (
		sb        strings.Builder // accumulates full answer
		eventType string          // current SSE event type
	)

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "data: [DONE]":
			// Stream finished successfully.
			return sb.String(), nil

		case strings.HasPrefix(line, "event:"):
			eventType = strings.TrimSpace(strings.TrimPrefix(line, "event:"))

		case strings.HasPrefix(line, "data:"):
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if err := s.handleSSEEvent(eventType, data, &sb, onDelta); err != nil {
				return "", err
			}

		case line == "":
			// Blank line resets event type for next event block.
			eventType = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("%w: stream read error: %v", ErrUpstreamFailed, err)
	}

	// Stream ended without [DONE] — return whatever was accumulated.
	answer := sb.String()
	if answer == "" {
		return "", fmt.Errorf("%w: empty answer from stream", ErrUpstreamFailed)
	}
	return answer, nil
}

// handleSSEEvent processes a single SSE event based on its type.
func (s *messageService) handleSSEEvent(eventType, data string, sb *strings.Builder, onDelta func(chunk string)) error {
	switch eventType {
	case "response.output_text.delta":
		var delta sseDeltaData
		if err := json.Unmarshal([]byte(data), &delta); err != nil {
			s.logger.Warn("failed to parse delta event", "data", data, "error", err)
			return nil // skip malformed delta, don't abort stream
		}
		if delta.Delta != "" {
			sb.WriteString(delta.Delta)
			if onDelta != nil {
				onDelta(delta.Delta)
			}
		}

	case "response.failed":
		return fmt.Errorf("%w: stream reported failure: %s", ErrUpstreamFailed, data)

	default:
		// Ignore lifecycle events (response.created, response.completed, etc.)
		s.logger.Debug("SSE event ignored", "type", eventType)
	}
	return nil
}

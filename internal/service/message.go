package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

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
	List(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error)
}

// openResponsesRequest is the body sent to POST /v1/responses.
type openResponsesRequest struct {
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
	Input  string `json:"input"`
}

type openResponsesMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openResponsesResp is the expected response shape from the OpenResponses API.
type openResponsesResp struct {
	Output []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

type messageService struct {
	cfg        *configs.Config
	repo       repository.MessageRepository
	httpClient *http.Client
	logger     *slog.Logger
}

// NewMessageService creates a MessageService with a 30-second HTTP timeout and structured logger.
func NewMessageService(cfg *configs.Config, repo repository.MessageRepository, logger *slog.Logger) MessageService {
	return &messageService{
		cfg:  cfg,
		repo: repo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Send calls the OpenResponses API with the user's question, saves the Q&A, and returns it.
func (s *messageService) Send(ctx context.Context, userID, question string) (*domain.Message, error) {
	if question == "" {
		return nil, ErrEmptyMessage
	}

	answer, err := s.callOpenResponses(ctx, question)
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

// callOpenResponses sends a POST /v1/responses request and extracts the answer text.
func (s *messageService) callOpenResponses(ctx context.Context, question string) (string, error) {
	payload := openResponsesRequest{
		Model:  s.cfg.ChatServer.Model,
		Stream: false,
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
		// Distinguish timeout from other network errors
		if errors.Is(err, context.DeadlineExceeded) {
			return "", ErrUpstreamTimeout
		}
		return "", fmt.Errorf("%w: %v", ErrUpstreamFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d", ErrUpstreamFailed, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("service.message: read response body: %w", err)
	}

	var parsed openResponsesResp
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("%w: failed to parse response JSON", ErrUpstreamFailed)
	}

	// Extract output[0].content[0].text
	if len(parsed.Output) == 0 || len(parsed.Output[0].Content) == 0 {
		return "", fmt.Errorf("%w: unexpected empty output", ErrUpstreamFailed)
	}
	answer := parsed.Output[0].Content[0].Text
	if answer == "" {
		return "", fmt.Errorf("%w: empty answer text in response", ErrUpstreamFailed)
	}

	return answer, nil
}

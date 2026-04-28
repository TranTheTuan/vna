package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/internal/domain"
	"github.com/TranTheTuan/vna/internal/dto"
	"github.com/TranTheTuan/vna/internal/service"
)

// mockThreadService is a test double for ThreadService.
type mockThreadService struct {
	createFn     func(ctx context.Context, userID string) (*domain.Thread, error)
	listByUserFn func(ctx context.Context, userID string) ([]*domain.Thread, error)
	renameFn     func(ctx context.Context, userID, threadID, name string) (*domain.Thread, error)
}

func (m *mockThreadService) Create(ctx context.Context, userID string) (*domain.Thread, error) {
	if m.createFn != nil {
		return m.createFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockThreadService) ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockThreadService) Rename(ctx context.Context, userID, threadID, name string) (*domain.Thread, error) {
	if m.renameFn != nil {
		return m.renameFn(ctx, userID, threadID, name)
	}
	return nil, nil
}

// mockMessageService is a test double for MessageService.
type mockMessageService struct {
	sendFn       func(ctx context.Context, userID, threadID, question string) (*domain.Message, error)
	sendStreamFn func(ctx context.Context, userID, threadID, question string, onMeta func(threadID string), onDelta func(chunk string)) (*domain.Message, error)
	listByThreadFn func(ctx context.Context, userID, threadID string, limit int, cursor string) ([]*domain.Message, string, error)
}

func (m *mockMessageService) Send(ctx context.Context, userID, threadID, question string) (*domain.Message, error) {
	if m.sendFn != nil {
		return m.sendFn(ctx, userID, threadID, question)
	}
	return nil, nil
}

func (m *mockMessageService) SendStream(ctx context.Context, userID, threadID, question string, onMeta func(threadID string), onDelta func(chunk string)) (*domain.Message, error) {
	if m.sendStreamFn != nil {
		return m.sendStreamFn(ctx, userID, threadID, question, onMeta, onDelta)
	}
	return nil, nil
}

func (m *mockMessageService) ListByThread(ctx context.Context, userID, threadID string, limit int, cursor string) ([]*domain.Message, string, error) {
	if m.listByThreadFn != nil {
		return m.listByThreadFn(ctx, userID, threadID, limit, cursor)
	}
	return nil, "", nil
}

func TestThreadHandler_List_Success(t *testing.T) {
	userID := "user123"
	expectedThreads := []*domain.Thread{
		{ID: "t1", UserID: userID, Name: "Chat 1", CreatedAt: time.Now()},
		{ID: "t2", UserID: userID, Name: "Chat 2", CreatedAt: time.Now()},
	}

	mockSvc := &mockThreadService{
		listByUserFn: func(ctx context.Context, uid string) ([]*domain.Thread, error) {
			return expectedThreads, nil
		},
	}

	handler := NewThreadHandler(mockSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	if err := handler.List(c); err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp dto.ListThreadsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Data) != len(expectedThreads) {
		t.Errorf("expected %d threads, got %d", len(expectedThreads), len(resp.Data))
	}

	for i, th := range resp.Data {
		if th.ID != expectedThreads[i].ID {
			t.Errorf("thread %d: expected ID %s, got %s", i, expectedThreads[i].ID, th.ID)
		}
	}
}

func TestThreadHandler_List_Empty(t *testing.T) {
	userID := "user123"

	mockSvc := &mockThreadService{
		listByUserFn: func(ctx context.Context, uid string) ([]*domain.Thread, error) {
			return []*domain.Thread{}, nil
		},
	}

	handler := NewThreadHandler(mockSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	if err := handler.List(c); err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp dto.ListThreadsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected 0 threads, got %d", len(resp.Data))
	}
}

func TestMessageHandler_ListByThread_Success(t *testing.T) {
	userID := "user123"
	threadID := "thread-abc"
	expectedMessages := []*domain.Message{
		{
			ID:        "msg1",
			UserID:    userID,
			ThreadID:  threadID,
			Question:  "Q1",
			Answer:    "A1",
			CreatedAt: time.Now(),
		},
	}

	mockSvc := &mockMessageService{
		listByThreadFn: func(ctx context.Context, uid, tid string, limit int, cursor string) ([]*domain.Message, string, error) {
			return expectedMessages, "", nil
		},
	}

	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/messages?thread_id="+threadID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	if err := handler.ListByThread(c); err != nil {
		t.Fatalf("ListByThread failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp dto.ListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Data) != len(expectedMessages) {
		t.Errorf("expected %d messages, got %d", len(expectedMessages), len(resp.Data))
	}
}

func TestMessageHandler_ListByThread_MissingThreadID(t *testing.T) {
	userID := "user123"

	mockSvc := &mockMessageService{}
	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	err := handler.ListByThread(c)
	if err == nil {
		t.Fatal("expected error for missing thread_id, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", httpErr.Code)
	}
}

func TestMessageHandler_ListByThread_ForeignThread(t *testing.T) {
	userID := "user123"
	threadID := "thread-foreign"

	mockSvc := &mockMessageService{
		listByThreadFn: func(ctx context.Context, uid, tid string, limit int, cursor string) ([]*domain.Message, string, error) {
			return nil, "", service.ErrThreadNotFound
		},
	}

	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/messages?thread_id="+threadID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	err := handler.ListByThread(c)
	if err == nil {
		t.Fatal("expected error for foreign thread, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", httpErr.Code)
	}
}

func TestMessageHandler_ListByThread_InvalidLimit(t *testing.T) {
	userID := "user123"
	threadID := "thread-abc"

	mockSvc := &mockMessageService{}
	handler := NewMessageHandler(mockSvc)

	tests := []struct {
		name  string
		limit string
	}{
		{"limit not a number", "abc"},
		{"limit too high", "101"},
		{"limit zero", "0"},
		{"limit negative", "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/messages?thread_id="+threadID+"&limit="+tt.limit, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", userID)

			err := handler.ListByThread(c)
			if err == nil {
				t.Fatal("expected error for invalid limit, got nil")
			}

			httpErr, ok := err.(*echo.HTTPError)
			if !ok {
				t.Fatalf("expected HTTPError, got %T", err)
			}

			if httpErr.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", httpErr.Code)
			}
		})
	}
}

func TestMessageHandler_ListByThread_WithCursor(t *testing.T) {
	userID := "user123"
	threadID := "thread-abc"
	cursor := "msg-last-id"
	limit := 10

	mockSvc := &mockMessageService{
		listByThreadFn: func(ctx context.Context, uid, tid string, lim int, cur string) ([]*domain.Message, string, error) {
			if cur != cursor {
				t.Errorf("expected cursor %s, got %s", cursor, cur)
			}
			if lim != limit {
				t.Errorf("expected limit %d, got %d", limit, lim)
			}
			return []*domain.Message{}, "", nil
		},
	}

	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	url := "/api/v1/messages?thread_id=" + threadID + "&limit=" + "10" + "&cursor=" + cursor
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	if err := handler.ListByThread(c); err != nil {
		t.Fatalf("ListByThread failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestMessageHandler_Send_Success(t *testing.T) {
	userID := "user123"
	threadID := "thread-abc"
	question := "What is Go?"

	mockSvc := &mockMessageService{
		sendFn: func(ctx context.Context, uid, tid, q string) (*domain.Message, error) {
			return &domain.Message{
				ID:        "msg1",
				UserID:    uid,
				ThreadID:  tid,
				Question:  q,
				Answer:    "Go is a programming language",
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"message":"What is Go?","thread_id":"thread-abc"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	if err := handler.Send(c); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}

	var resp dto.MessageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Question != question {
		t.Errorf("expected question %s, got %s", question, resp.Question)
	}
	if resp.ThreadID != threadID {
		t.Errorf("expected thread ID %s, got %s", threadID, resp.ThreadID)
	}
}

func TestMessageHandler_Send_EmptyMessage(t *testing.T) {
	userID := "user123"

	mockSvc := &mockMessageService{}
	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"message":"","thread_id":"thread-abc"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	err := handler.Send(c)
	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", httpErr.Code)
	}
}

func TestMessageHandler_Send_ThreadNotFound(t *testing.T) {
	userID := "user123"

	mockSvc := &mockMessageService{
		sendFn: func(ctx context.Context, uid, tid, q string) (*domain.Message, error) {
			return nil, service.ErrThreadNotFound
		},
	}

	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"message":"Hello","thread_id":"thread-foreign"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	err := handler.Send(c)
	if err == nil {
		t.Fatal("expected error for thread not found, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", httpErr.Code)
	}
}

func TestMessageHandler_SendStream_MetadataFirst(t *testing.T) {
	userID := "user123"
	threadID := "thread-abc"

	metaEmitted := false
	var metaThreadID string

	mockSvc := &mockMessageService{
		sendStreamFn: func(ctx context.Context, uid, tid, q string, onMeta func(string), onDelta func(string)) (*domain.Message, error) {
			// onMeta should be called first, before any onDelta
			if onMeta != nil {
				onMeta(tid)
				metaEmitted = true
				metaThreadID = tid
			}
			// Simulate some delta events
			if onDelta != nil {
				onDelta("Hello ")
				onDelta("world")
			}
			return &domain.Message{
				ID:        "msg1",
				UserID:    uid,
				ThreadID:  tid,
				Question:  q,
				Answer:    "Hello world",
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"message":"What is Go?","thread_id":"thread-abc"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages/stream", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	if err := handler.SendStream(c); err != nil {
		t.Fatalf("SendStream failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if !metaEmitted {
		t.Fatal("expected onMeta to be called")
	}

	if metaThreadID != threadID {
		t.Errorf("expected thread ID %s in metadata, got %s", threadID, metaThreadID)
	}

	// Check that response contains SSE events
	responseBody := rec.Body.String()
	if !strings.Contains(responseBody, "event: metadata") {
		t.Error("expected 'event: metadata' in response")
	}
	if !strings.Contains(responseBody, "event: delta") {
		t.Error("expected 'event: delta' in response")
	}
	if !strings.Contains(responseBody, "event: done") {
		t.Error("expected 'event: done' in response")
	}
}

func TestMessageHandler_SendStream_EmptyMessage(t *testing.T) {
	userID := "user123"

	mockSvc := &mockMessageService{}
	handler := NewMessageHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"message":"","thread_id":"thread-abc"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages/stream", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)

	err := handler.SendStream(c)
	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", httpErr.Code)
	}
}

func TestThreadHandler_Rename_Success(t *testing.T) {
	userID := "user123"
	threadID := "thread-abc"
	newName := "Updated Chat"

	mockSvc := &mockThreadService{
		renameFn: func(ctx context.Context, uid, tid, name string) (*domain.Thread, error) {
			if uid != userID {
				t.Errorf("expected userID %s, got %s", userID, uid)
			}
			if tid != threadID {
				t.Errorf("expected threadID %s, got %s", threadID, tid)
			}
			if name != newName {
				t.Errorf("expected name %s, got %s", newName, name)
			}
			return &domain.Thread{
				ID:        tid,
				UserID:    uid,
				Name:      name,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewThreadHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"name":"Updated Chat"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/threads/"+threadID, body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)
	c.SetParamNames("id")
	c.SetParamValues(threadID)

	if err := handler.Rename(c); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp dto.ThreadResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.ID != threadID {
		t.Errorf("expected thread ID %s, got %s", threadID, resp.ID)
	}
	if resp.Name != newName {
		t.Errorf("expected name %s, got %s", newName, resp.Name)
	}
}

func TestThreadHandler_Rename_EmptyName(t *testing.T) {
	userID := "user123"
	threadID := "thread-abc"

	mockSvc := &mockThreadService{}
	handler := NewThreadHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"name":""}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/threads/"+threadID, body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)
	c.SetParamNames("id")
	c.SetParamValues(threadID)

	err := handler.Rename(c)
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", httpErr.Code)
	}
}

func TestThreadHandler_Rename_ThreadNotFound(t *testing.T) {
	userID := "user123"
	threadID := "thread-foreign"

	mockSvc := &mockThreadService{
		renameFn: func(ctx context.Context, uid, tid, name string) (*domain.Thread, error) {
			return nil, service.ErrThreadNotFound
		},
	}

	handler := NewThreadHandler(mockSvc)

	e := echo.New()
	body := strings.NewReader(`{"name":"New Name"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/threads/"+threadID, body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)
	c.SetParamNames("id")
	c.SetParamValues(threadID)

	err := handler.Rename(c)
	if err == nil {
		t.Fatal("expected error for thread not found, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", httpErr.Code)
	}
}

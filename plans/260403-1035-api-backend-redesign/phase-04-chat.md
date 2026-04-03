# Phase 04 — Chat (Send Message + OpenResponses)

**Priority:** P1 | **Status:** Pending | **Effort:** 1.5h

## Overview

Implement the message-send endpoint: receive user message → call OpenResponses `POST /v1/responses` → save question+answer → return answer.

## Context Links

- Brainstorm: `plans/reports/brainstorm-260403-0933-api-backend-redesign.md`
- Depends on Phase 02 (domain/Message, migration 002) and Phase 03 (JWT middleware exists)

## API Contract

```
POST /api/v1/messages
  Header:  Authorization: Bearer <access_token>
  Body:    { "message": "What is the capital of France?" }
  Success: 201 {
    "id": "uuid",
    "question": "What is the capital of France?",
    "answer": "Paris.",
    "created_at": "2026-04-03T10:00:00Z"
  }
  Errors:
    400 missing/empty message
    401 unauthorized
    502 OpenResponses API error
    504 OpenResponses timeout
```

## OpenResponses API Spec

```
POST {OPENRESPONSES_URL}/v1/responses
Headers:
  Authorization: Bearer {OPENRESPONSES_API_KEY}
  Content-Type: application/json
Body:
  {
    "model": "{OPENRESPONSES_MODEL}",
    "input": [{ "role": "user", "content": "<user message>" }]
  }
Response (200 full JSON):
  {
    "output": [
      {
        "content": [
          { "type": "output_text", "text": "<answer string>" }
        ]
      }
    ]
  }
```

Extract path: `response.output[0].content[0].text`

HTTP client: `net/http` stdlib. Timeout: 30s (set on `http.Client`).

## Files to IMPLEMENT

### `internal/dto/message_dto.go`

```go
type SendMessageRequest  { Message string `json:"message"` }
type MessageResponse     {
    ID        string    `json:"id"`
    Question  string    `json:"question"`
    Answer    string    `json:"answer"`
    CreatedAt time.Time `json:"created_at"`
}
```

### `internal/repository/message.go`

```go
type MessageRepository interface {
    Save(ctx context.Context, msg *domain.Message) (*domain.Message, error)
    ListByUser(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error)
}
```

SQL for Save:
```sql
INSERT INTO messages(user_id, question, answer)
VALUES($1, $2, $3)
RETURNING id, user_id, question, answer, created_at
```

`ListByUser` implemented in Phase 05.

### `internal/service/message.go`

```go
type MessageService interface {
    Send(ctx context.Context, userID, question string) (*domain.Message, error)
    List(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error)
}
```

`Send` implementation:
1. Validate question not empty
2. Build OpenResponses request struct
3. Marshal to JSON, POST to `{cfg.OpenResponsesURL}/v1/responses` with 30s timeout
4. Parse response, extract `output[0].content[0].text`
5. If HTTP status != 200 → return `ErrUpstreamFailed` (caller maps to 502)
6. If timeout → `ErrUpstreamTimeout` (caller maps to 504)
7. Call `repo.Save(ctx, &domain.Message{UserID: userID, Question: question, Answer: answer})`
8. Return saved message

Define response structs internally (private) for JSON parsing:
```go
type openResponsesResp struct {
    Output []struct {
        Content []struct {
            Type string `json:"type"`
            Text string `json:"text"`
        } `json:"content"`
    } `json:"output"`
}
```

Constructor: `NewMessageService(cfg *configs.Config, repo repository.MessageRepository) *messageService`
Inject `http.Client` with 30s timeout.

### `internal/handler/http/message.go`

```go
type MessageHandler struct { svc service.MessageService }

func (h *MessageHandler) Send(c echo.Context) error   // POST /api/v1/messages
func (h *MessageHandler) List(c echo.Context) error   // GET /api/v1/messages (Phase 05)
```

`Send`:
1. Bind `SendMessageRequest`, validate `Message` not empty
2. Get `user_id` from context (set by JWT middleware)
3. Call `svc.Send`
4. Return 201 `MessageResponse`
5. Map `ErrUpstreamFailed` → 502, `ErrUpstreamTimeout` → 504

## Todo List

- [ ] Create `internal/dto/message_dto.go`
- [ ] Implement `internal/repository/message.go` (Save method; ListByUser stub)
- [ ] Implement `internal/service/message.go` (Send only; List stub for Phase 05)
- [ ] Create `internal/handler/http/message.go` (Send handler; List stub)
- [ ] `go build ./internal/...` compiles cleanly

## Success Criteria

- Handler calls OpenResponses with correct payload
- On 200 response: message saved to DB, answer returned
- On upstream error (non-200): 502 returned to client, nothing saved
- On timeout: 504 returned, nothing saved
- Empty message body: 400 returned immediately

## Risk Assessment

| Risk | Mitigation |
|---|---|
| OpenResponses schema change (output path) | Parse defensively; return 502 if path missing |
| API key leaked in logs | Never log request body or Authorization header |
| Large answers | TEXT column handles unlimited; no truncation |

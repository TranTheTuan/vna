# Brainstorm Report: Multi-Thread Chat

**Date:** 2026-04-21  
**Status:** Agreed

---

## Problem Statement

Currently one flat chat history per user. Users want multiple named threads (like ChatGPT's conversation sidebar). Need to scope messages, OpenResponses API calls, and pagination per thread.

---

## Evaluated Approaches

### Option A — Implicit via empty thread_id (user's original proposal)
- `POST /api/v1/messages/stream` with `thread_id: ""` auto-creates thread
- `thread_id` returned in `done` event
- **Con:** Thread ID arrives late (after streaming ends); no sidebar entry until first message completes; race condition if two messages sent simultaneously

### Option B — Explicit POST /threads
- Dedicated create-thread endpoint; client gets `thread_id` before first message
- **Con:** Extra round-trip

### Option C — Hybrid (chosen with refinement)
- Empty `thread_id` = create new; existing UUID = use existing
- Thread ID streamed back in **`event: metadata`** as the FIRST SSE event — before any delta
- Client gets thread_id immediately, can update sidebar without waiting for stream to finish

---

## Final Recommended Solution

### API Contract

```
POST /api/v1/messages/stream
Body: { "message": "...", "thread_id": "" | "<uuid>" }

SSE output order:
  event: metadata
  data: {"thread_id":"<uuid>"}            ← always first

  event: delta
  data: {"delta":"<chunk>"}               ← 0..N

  event: done
  data: {"id":"...","question":"...","answer":"...","thread_id":"...","created_at":"..."}

  event: error                            ← replaces done on failure
  data: {"message":"..."}

GET /api/v1/threads
Headers: Authorization: Bearer <token>
Response: {"data":[{"id":"...","name":"New Chat","created_at":"..."}]}

GET /api/v1/messages?thread_id=<uuid>&limit=20&cursor=<uuid>
Response: {"data":[...],"next_cursor":"..."}
```

### Database Changes

```sql
-- New table
CREATE TABLE threads (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL DEFAULT 'New Chat',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_threads_user ON threads(user_id, created_at DESC);

-- Extend messages (nullable for backward compat)
ALTER TABLE messages ADD COLUMN thread_id UUID REFERENCES threads(id) ON DELETE CASCADE;
CREATE INDEX idx_messages_thread_time ON messages(thread_id, created_at DESC);
```

### Server Flow (SendStream handler)

```
1. Parse {message, thread_id} from body
2. If thread_id == "":
     INSERT INTO threads(user_id, name) VALUES($1, 'New Chat') RETURNING id
   Else:
     SELECT 1 FROM threads WHERE id=$1 AND user_id=$2  → 404 if not found
3. Emit: event: metadata  data: {"thread_id":"<id>"}
4. Call OpenResponses API with user=thread_id (stream=true)
5. Emit: event: delta per chunk
6. INSERT INTO messages(user_id, thread_id, question, answer) RETURNING *
7. Emit: event: done with full MessageResponse (including thread_id)
```

### Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Creation trigger | Empty `thread_id` on stream call | No extra round-trip |
| Thread ID delivery | `event: metadata` before first delta | Client can update sidebar immediately |
| Thread name | Server default "New Chat" | KISS; rename via PATCH later if needed |
| OpenResponses `user` param | `thread_id` | Per-thread AI context isolation |
| `GET /messages` scope | `?thread_id` required | Messages always thread-scoped |
| Legacy messages | `thread_id` nullable | Safe migration |

---

## Implementation Considerations

- `SendMessageRequest` DTO gains optional `ThreadID string \`json:"thread_id"\`` field
- `MessageResponse` DTO gains `ThreadID string \`json:"thread_id"\``
- New `ThreadRepository` + `ThreadService` interfaces (follows existing layered pattern)
- New migration file: `003_threads.sql`
- `GET /api/v1/messages` changes: `thread_id` becomes a required query param (breaking change — coordinate with frontend)
- Validate `thread_id` ownership before use to prevent cross-user thread access

## Risks

- Breaking change on `GET /api/v1/messages` — must update frontend simultaneously
- OpenResponses `user` param change (userID → threadID): ensure the upstream API doesn't confuse existing user context

## Next Steps

- Create implementation plan via `/plan`
- New files: `domain/thread.go`, `repository/thread.go`, `service/thread.go`, `handler/http/thread.go`, `dto/thread_dto.go`, `delivery/http/thread_routes.go`
- New migration: `internal/migrations/003_threads.sql`
- Update: `service/message.go`, `repository/message.go`, `handler/http/message.go`, `dto/message_dto.go`

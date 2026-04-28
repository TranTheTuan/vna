# Phase 01 — Database Migration

**Context:** `plans/reports/brainstorm-260421-0848-multi-thread-chat.md`
**Priority:** P1 | **Status:** Pending | **Effort:** 0.5h

## Overview

Add `threads` table and `thread_id` FK column to `messages`. Nullable FK preserves backward compat for existing rows.

## Requirements

- `threads(id, user_id, name, created_at)` with index on `(user_id, created_at DESC)`
- `messages.thread_id UUID NULL REFERENCES threads(id) ON DELETE CASCADE`
- Index on `messages(thread_id, created_at DESC)` for keyset pagination

## Related Code Files

| Action | File |
|--------|------|
| Create | `internal/migrations/003_threads.up.sql` |

## Implementation Steps

1. Create `internal/migrations/003_threads.up.sql`:

```sql
-- Migration 003: threads table and thread_id on messages

CREATE TABLE IF NOT EXISTS threads (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL DEFAULT 'New Chat',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_threads_user
    ON threads(user_id, created_at DESC);

ALTER TABLE messages ADD COLUMN IF NOT EXISTS thread_id UUID REFERENCES threads(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_messages_thread_time
    ON messages(thread_id, created_at DESC);
```

2. Apply migration (manually or via `make migrate` if wired):
   ```bash
   psql $DATABASE_URL -f internal/migrations/003_threads.up.sql
   ```

## Todo

- [x] Create `internal/migrations/003_threads.up.sql`
- [x] Apply migration to dev DB and verify schema

## Success Criteria

- `threads` table exists with correct columns and indexes
- `messages.thread_id` column exists (nullable)
- Existing message rows unaffected (`thread_id = NULL`)

## Risk Assessment

- `ADD COLUMN IF NOT EXISTS` is safe on existing data
- `NULL` default avoids locking issues

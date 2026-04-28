package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/TranTheTuan/vna/internal/domain"
)

// mockThreadRepository is a test double for ThreadRepository.
type mockThreadRepository struct {
	createFn         func(ctx context.Context, userID string) (*domain.Thread, error)
	listByUserFn     func(ctx context.Context, userID string) ([]*domain.Thread, error)
	getByIDAndUserFn func(ctx context.Context, threadID, userID string) (*domain.Thread, error)
	deleteFn         func(ctx context.Context, threadID string) error
	renameFn         func(ctx context.Context, threadID, name string) (*domain.Thread, error)
}

func (m *mockThreadRepository) Create(ctx context.Context, userID string) (*domain.Thread, error) {
	if m.createFn != nil {
		return m.createFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockThreadRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockThreadRepository) GetByIDAndUser(ctx context.Context, threadID, userID string) (*domain.Thread, error) {
	if m.getByIDAndUserFn != nil {
		return m.getByIDAndUserFn(ctx, threadID, userID)
	}
	return nil, nil
}

func (m *mockThreadRepository) Delete(ctx context.Context, threadID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, threadID)
	}
	return nil
}

func (m *mockThreadRepository) Rename(ctx context.Context, threadID, name string) (*domain.Thread, error) {
	if m.renameFn != nil {
		return m.renameFn(ctx, threadID, name)
	}
	return nil, nil
}

func TestThreadService_Create(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	threadID := "thread-abc"

	mock := &mockThreadRepository{
		createFn: func(ctx context.Context, uid string) (*domain.Thread, error) {
			if uid != userID {
				t.Errorf("expected userID %s, got %s", userID, uid)
			}
			return &domain.Thread{
				ID:        threadID,
				UserID:    uid,
				Name:      "New Chat",
				CreatedAt: time.Now(),
			}, nil
		},
	}

	svc := NewThreadService(mock)
	thread, err := svc.Create(ctx, userID)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if thread.ID != threadID {
		t.Errorf("expected thread ID %s, got %s", threadID, thread.ID)
	}
	if thread.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, thread.UserID)
	}
}

func TestThreadService_ListByUser(t *testing.T) {
	ctx := context.Background()
	userID := "user123"

	expectedThreads := []*domain.Thread{
		{ID: "t1", UserID: userID, Name: "Chat 1", CreatedAt: time.Now()},
		{ID: "t2", UserID: userID, Name: "Chat 2", CreatedAt: time.Now()},
	}

	mock := &mockThreadRepository{
		listByUserFn: func(ctx context.Context, uid string) ([]*domain.Thread, error) {
			if uid != userID {
				t.Errorf("expected userID %s, got %s", userID, uid)
			}
			return expectedThreads, nil
		},
	}

	svc := NewThreadService(mock)
	threads, err := svc.ListByUser(ctx, userID)

	if err != nil {
		t.Fatalf("ListByUser failed: %v", err)
	}
	if len(threads) != len(expectedThreads) {
		t.Errorf("expected %d threads, got %d", len(expectedThreads), len(threads))
	}
	for i, th := range threads {
		if th.ID != expectedThreads[i].ID {
			t.Errorf("thread %d: expected ID %s, got %s", i, expectedThreads[i].ID, th.ID)
		}
	}
}

func TestThreadService_ListByUser_Empty(t *testing.T) {
	ctx := context.Background()
	userID := "user123"

	mock := &mockThreadRepository{
		listByUserFn: func(ctx context.Context, uid string) ([]*domain.Thread, error) {
			return []*domain.Thread{}, nil
		},
	}

	svc := NewThreadService(mock)
	threads, err := svc.ListByUser(ctx, userID)

	if err != nil {
		t.Fatalf("ListByUser failed: %v", err)
	}
	if len(threads) != 0 {
		t.Errorf("expected 0 threads, got %d", len(threads))
	}
}

func TestThreadService_Rename_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	threadID := "thread-abc"
	newName := "Updated Chat"

	mock := &mockThreadRepository{
		getByIDAndUserFn: func(ctx context.Context, tid, uid string) (*domain.Thread, error) {
			if tid != threadID {
				t.Errorf("expected threadID %s, got %s", threadID, tid)
			}
			if uid != userID {
				t.Errorf("expected userID %s, got %s", userID, uid)
			}
			return &domain.Thread{
				ID:        tid,
				UserID:    uid,
				Name:      "Old Name",
				CreatedAt: time.Now(),
			}, nil
		},
		renameFn: func(ctx context.Context, tid, name string) (*domain.Thread, error) {
			if tid != threadID {
				t.Errorf("expected threadID %s, got %s", threadID, tid)
			}
			if name != newName {
				t.Errorf("expected name %s, got %s", newName, name)
			}
			return &domain.Thread{
				ID:        tid,
				UserID:    userID,
				Name:      name,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	svc := NewThreadService(mock)
	thread, err := svc.Rename(ctx, userID, threadID, newName)

	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}
	if thread.ID != threadID {
		t.Errorf("expected thread ID %s, got %s", threadID, thread.ID)
	}
	if thread.Name != newName {
		t.Errorf("expected name %s, got %s", newName, thread.Name)
	}
}

func TestThreadService_Rename_InvalidName(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	threadID := "thread-abc"

	mock := &mockThreadRepository{}
	svc := NewThreadService(mock)

	_, err := svc.Rename(ctx, userID, threadID, "")

	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
	if err != ErrInvalidThreadName {
		t.Errorf("expected ErrInvalidThreadName, got %v", err)
	}
}

func TestThreadService_Rename_ForeignThread(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	threadID := "thread-foreign"

	mock := &mockThreadRepository{
		getByIDAndUserFn: func(ctx context.Context, tid, uid string) (*domain.Thread, error) {
			// Simulate thread not found for this user
			return nil, sql.ErrNoRows
		},
	}

	svc := NewThreadService(mock)
	_, err := svc.Rename(ctx, userID, threadID, "New Name")

	if err == nil {
		t.Fatal("expected error for foreign thread, got nil")
	}
	if err != ErrThreadNotFound {
		t.Errorf("expected ErrThreadNotFound, got %v", err)
	}
}

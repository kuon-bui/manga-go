package authorizationadmin

import (
	"context"
	"database/sql"
	"sync"

	"manga-go/internal/app/api/common/response"

	"gorm.io/gorm"
)

const authorizationMutationLockKey int64 = 0x415554485A

type MutationLocker interface {
	WithLock(ctx context.Context, fn func() response.Result) response.Result
}

type postgresMutationLocker struct {
	db *gorm.DB
}

func NewPostgresMutationLocker(db *gorm.DB) MutationLocker {
	return &postgresMutationLocker{db: db}
}

func (l *postgresMutationLocker) WithLock(ctx context.Context, fn func() response.Result) response.Result {
	if l == nil || l.db == nil {
		return response.ResultErrInternal(gorm.ErrInvalidDB)
	}
	db, err := l.db.DB()
	if err != nil {
		return response.ResultErrInternal(err)
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return response.ResultErrInternal(err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", authorizationMutationLockKey); err != nil {
		return response.ResultErrInternal(err)
	}
	defer releaseAdvisoryLock(conn)
	return fn()
}

func releaseAdvisoryLock(conn *sql.Conn) {
	_, _ = conn.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", authorizationMutationLockKey)
}

type mutexMutationLocker struct {
	mu sync.Mutex
}

func NewMutexMutationLocker() MutationLocker {
	return &mutexMutationLocker{}
}

func (l *mutexMutationLocker) WithLock(_ context.Context, fn func() response.Result) response.Result {
	l.mu.Lock()
	defer l.mu.Unlock()
	return fn()
}

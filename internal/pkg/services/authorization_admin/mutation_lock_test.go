package authorizationadmin

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"manga-go/internal/app/api/common/response"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type advisoryLockExec struct {
	connectionID int64
	query        string
	key          int64
}

func TestPostgresMutationLockerSerializesIndependentConnections(t *testing.T) {
	dsn := os.Getenv("AUTHORIZATION_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("set AUTHORIZATION_TEST_POSTGRES_DSN to run PostgreSQL advisory-lock contention coverage")
	}
	firstDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	firstSQLDB, err := firstDB.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = firstSQLDB.Close() })
	secondDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	secondSQLDB, err := secondDB.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = secondSQLDB.Close() })
	firstEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondAttempting := make(chan struct{})
	secondEntered := make(chan struct{})
	firstDone := make(chan response.Result, 1)
	secondDone := make(chan response.Result, 1)

	go func() {
		firstDone <- NewPostgresMutationLocker(firstDB).WithLock(context.Background(), func() response.Result {
			close(firstEntered)
			<-releaseFirst
			return response.ResultSuccess("first", nil)
		})
	}()
	select {
	case <-firstEntered:
	case result := <-firstDone:
		t.Fatalf("first PostgreSQL locker failed before entering: %#v", result)
	case <-time.After(5 * time.Second):
		t.Fatal("first PostgreSQL locker did not enter")
	}
	go func() {
		close(secondAttempting)
		secondDone <- NewPostgresMutationLocker(secondDB).WithLock(context.Background(), func() response.Result {
			close(secondEntered)
			return response.ResultSuccess("second", nil)
		})
	}()
	<-secondAttempting

	select {
	case <-secondEntered:
		t.Fatal("second independent PostgreSQL connection entered while the advisory lock was held")
	case result := <-secondDone:
		t.Fatalf("second PostgreSQL locker finished while the advisory lock was held: %#v", result)
	case <-time.After(100 * time.Millisecond):
	}
	close(releaseFirst)
	select {
	case <-secondEntered:
	case result := <-secondDone:
		t.Fatalf("second PostgreSQL locker failed before entering: %#v", result)
	case <-time.After(5 * time.Second):
		t.Fatal("second PostgreSQL connection did not enter after advisory lock release")
	}
	if result := <-firstDone; !result.Success {
		t.Fatalf("unexpected first advisory-lock result: %#v", result)
	}
	if result := <-secondDone; !result.Success {
		t.Fatalf("unexpected second advisory-lock result: %#v", result)
	}
}

type advisoryLockDriverState struct {
	mu     sync.Mutex
	nextID int64
	execs  []advisoryLockExec
}

func (s *advisoryLockDriverState) append(item advisoryLockExec) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.execs = append(s.execs, item)
}

type advisoryLockDriver struct {
	state *advisoryLockDriverState
}

func (d *advisoryLockDriver) Open(string) (driver.Conn, error) {
	id := atomic.AddInt64(&d.state.nextID, 1)
	return &advisoryLockConn{id: id, state: d.state}, nil
}

type advisoryLockConn struct {
	id    int64
	state *advisoryLockDriverState
}

func (c *advisoryLockConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not supported")
}

func (c *advisoryLockConn) Close() error { return nil }

func (c *advisoryLockConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are not supported")
}

func (c *advisoryLockConn) ExecContext(
	_ context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Result, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected one advisory lock key, got %d", len(args))
	}
	key, ok := args[0].Value.(int64)
	if !ok {
		return nil, fmt.Errorf("expected int64 advisory lock key, got %T", args[0].Value)
	}
	c.state.append(advisoryLockExec{connectionID: c.id, query: query, key: key})
	return driver.RowsAffected(1), nil
}

func TestPostgresMutationLockerUsesOneDedicatedSessionForLockAndUnlock(t *testing.T) {
	state := &advisoryLockDriverState{}
	driverName := fmt.Sprintf("authorization-advisory-lock-%d", atomic.AddInt64(&state.nextID, 1))
	sql.Register(driverName, &advisoryLockDriver{state: state})
	sqlDB, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB}),
		&gorm.Config{DisableAutomaticPing: true},
	)
	if err != nil {
		t.Fatal(err)
	}

	called := false
	result := NewPostgresMutationLocker(db).WithLock(context.Background(), func() response.Result {
		called = true
		return response.ResultSuccess("done", nil)
	})
	if !result.Success || !called {
		t.Fatalf("expected protected callback success, got %#v", result)
	}
	if len(state.execs) != 2 {
		t.Fatalf("expected advisory lock and unlock, got %#v", state.execs)
	}
	if state.execs[0].query != "SELECT pg_advisory_lock($1)" ||
		state.execs[1].query != "SELECT pg_advisory_unlock($1)" {
		t.Fatalf("unexpected advisory lock SQL: %#v", state.execs)
	}
	if state.execs[0].connectionID != state.execs[1].connectionID {
		t.Fatalf("lock and unlock used different sessions: %#v", state.execs)
	}
	if state.execs[0].key != authorizationMutationLockKey || state.execs[1].key != authorizationMutationLockKey {
		t.Fatalf("unexpected advisory lock keys: %#v", state.execs)
	}
}

var _ driver.ExecerContext = (*advisoryLockConn)(nil)

package queue_test

import (
	"testing"

	"manga-go/internal/queue"
	queueasynq "manga-go/internal/queue/asynq"
)

func TestQueueModulesAreRegistered(t *testing.T) {
	if queue.Module == nil {
		t.Fatal("expected queue.Module to be non-nil")
	}
	if queueasynq.Module == nil {
		t.Fatal("expected queue/asynq.Module to be non-nil")
	}
}

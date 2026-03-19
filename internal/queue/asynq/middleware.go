package asynq

import (
	"context"
	"fmt"
	r "reflect"

	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func tracerMiddleware(handler asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		taskType := task.Type()
		queueName := r.Indirect(r.ValueOf(task.ResultWriter())).FieldByName("qname").String()
		spanName := fmt.Sprintf("asynq %s %s", queueName, taskType)
		ctxWithTracer, span := otel.Tracer("asynq").Start(ctx, spanName)
		defer span.End()

		span.SetAttributes(attribute.String("task.queue", queueName))
		span.SetAttributes(attribute.String("task.type", taskType))

		return handler.ProcessTask(ctxWithTracer, task)
	})
}

package response

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestParseValidationErrorsNil(t *testing.T) {
	if got := parseValidationErrors(nil); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

func TestParseValidationErrorsNonValidatorError(t *testing.T) {
	if got := parseValidationErrors(errors.New("plain error")); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

func TestParseValidationErrorsRequired(t *testing.T) {
	type payload struct {
		Email string `validate:"required"`
	}

	v := validator.New()
	err := v.Struct(payload{})

	got := parseValidationErrors(err)
	if len(got) != 1 {
		t.Fatalf("expected 1 validation error, got %d", len(got))
	}
	if got[0].Field != "Email" {
		t.Fatalf("expected field Email, got %s", got[0].Field)
	}
	if got[0].Message != "Email is required" {
		t.Fatalf("expected required message, got %s", got[0].Message)
	}
}

func TestParseValidationErrorsUUID(t *testing.T) {
	type payload struct {
		ID string `validate:"uuid4"`
	}

	v := validator.New()
	err := v.Struct(payload{ID: "not-a-uuid"})

	got := parseValidationErrors(err)
	if len(got) != 1 {
		t.Fatalf("expected 1 validation error, got %d", len(got))
	}
	if got[0].Field != "ID" {
		t.Fatalf("expected field ID, got %s", got[0].Field)
	}
	if got[0].Message != "ID must be a valid UUID" {
		t.Fatalf("expected uuid message, got %s", got[0].Message)
	}
}

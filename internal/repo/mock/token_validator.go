package mock

import (
	"context"

	"github.com/pdaccess/pvault/internal/core/ports"
)

type mockTokenValidator struct{}

func NewAllValidValidator() ports.TokenValidator {
	return &mockTokenValidator{}
}

func (*mockTokenValidator) Validate(ctx context.Context, token string) error {
	return nil
}

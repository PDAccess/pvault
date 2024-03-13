package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/pdaccess/pvault/internal/core/domain"
)

func (s *ServiceImpl) auth(ctx context.Context) error {
	tokenValue, ok := ctx.Value(domain.UserTokenIn).(string)

	if !ok || strings.TrimSpace(tokenValue) == "" {
		return fmt.Errorf("no token in context: %w", domain.ErrNotFound)
	}

	for _, tokenValidator := range s.validators {
		return tokenValidator.Validate(ctx, tokenValue)
	}

	return domain.ErrInvalidToken
}

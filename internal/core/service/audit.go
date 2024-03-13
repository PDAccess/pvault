package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/pdaccess/pvault/internal/core/domain"
)

// Logs implements ports.Service.
func (s *ServiceImpl) Logs(ctx context.Context, key domain.RecordKey) ([]domain.RawAudit, error) {

	if err := s.auth(ctx); err != nil {
		return nil, fmt.Errorf("unauthenticated: %w", err)
	}

	aduits, err := s.persistence.SearchAudit(ctx, string(key))
	if err != nil {
		return nil, fmt.Errorf("logs: %w", err)
	}

	var ret []domain.RawAudit

	for _, audit := range aduits {
		decrtyped, err := s.decrypt(audit)

		if err != nil {
			return nil, fmt.Errorf("decrypt: %w", err)
		}

		var d domain.RawAudit
		if err := d.Restore(decrtyped); err != nil {
			return nil, fmt.Errorf("restore: %w", err)
		}

		ret = append(ret, d)
	}

	slices.SortFunc(ret, domain.AuditSortByOpTime)

	return ret, nil
}

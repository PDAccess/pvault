package ports

import (
	"context"

	"github.com/pdaccess/pvault/internal/core/domain"
)

type Service interface {
	Write(ctx context.Context, key domain.RecordKey, record domain.SimpleRecord) error
	Read(ctx context.Context, key domain.RecordKey) (*domain.SimpleRecord, error)

	Logs(ctx context.Context, key domain.RecordKey) ([]domain.RawAudit, error)
}

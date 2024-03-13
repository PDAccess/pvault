package ports

import (
	"context"
)

type Persistance interface {
	UpsertRecord(ctx context.Context, parent, key string, value []byte) error
	AppendAudit(ctx context.Context, key string, auditData []byte) error

	FetchRecord(ctx context.Context, parent, key string) ([]byte, error)
	SearchAudit(ctx context.Context, key string) ([][]byte, error)
}

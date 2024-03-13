package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pdaccess/pvault/internal/core/domain"
)

// Read implements ports.Service.
func (s *ServiceImpl) Read(ctx context.Context, key domain.RecordKey) (*domain.SimpleRecord, error) {

	if err := s.auth(ctx); err != nil {
		return nil, fmt.Errorf("unauthenticated: %w", err)
	}

	parent, keyStr := key.Split()
	encryted, err := s.persistence.FetchRecord(ctx, parent, keyStr)

	if err != nil {
		return nil, fmt.Errorf("FetchRecord: %w", err)
	}

	decrypt, err := s.decrypt(encryted)

	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	rec := &domain.SimpleRecord{}

	if err := rec.Restore(decrypt); err != nil {
		return nil, fmt.Errorf("restore: %w", err)
	}

	if err := s.trackAudit(ctx, rec.Id, domain.OpTypeRead); err != nil {
		return nil, fmt.Errorf("trackAudit: %w", err)
	}

	return rec, nil
}

// Write implements ports.Service.
func (s *ServiceImpl) Write(ctx context.Context, key domain.RecordKey, record domain.SimpleRecord) error {
	if err := s.auth(ctx); err != nil {
		return fmt.Errorf("unauthenticated: %w", err)
	}

	parent, keyStr := key.Split()

	bytes, err := record.Bytes()

	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	encrypted, err := s.encrypt(bytes)

	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	if err := s.persistence.UpsertRecord(ctx, parent, keyStr, encrypted); err != nil {
		return fmt.Errorf("upsert record: %w", err)
	}

	if err := s.trackAudit(ctx, keyStr, domain.OpTypeWrite); err != nil {
		return fmt.Errorf("trackAudit: %w", err)
	}

	return nil
}

func (s *ServiceImpl) trackAudit(ctx context.Context, keyStr string, opType domain.OpType) error {
	audit := domain.RawAudit{
		OperationTime: time.Now(),
		OpType:        opType,
	}

	bytes, err := audit.Bytes()

	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	encrypted, err := s.encrypt(bytes)

	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	if err := s.persistence.AppendAudit(ctx, keyStr, encrypted); err != nil {
		return fmt.Errorf("append audit: %w", err)
	}

	return nil
}

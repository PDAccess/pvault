package mock

import (
	"context"
	"fmt"

	"github.com/pdaccess/pvault/internal/core/domain"
	"github.com/pdaccess/pvault/internal/core/ports"
)

type mockDb struct {
	records map[string]map[string][]byte
	audits  map[string][][]byte
}

// AppendAudit implements ports.Persistance.
func (m *mockDb) AppendAudit(ctx context.Context, key string, auditData []byte) error {
	m.audits[key] = append(m.audits[key], auditData)

	return nil
}

// FetchRecord implements ports.Persistance.
func (m *mockDb) FetchRecord(ctx context.Context, parent string, key string) ([]byte, error) {
	p, ok := m.records[parent]

	if !ok {
		return nil, fmt.Errorf("not found: %w", domain.ErrNotFound)
	}

	return p[key], nil
}

// SearchAudit implements ports.Persistance.
func (m *mockDb) SearchAudit(ctx context.Context, key string) ([][]byte, error) {
	return m.audits[key], nil
}

// UpsertRecord implements ports.Persistance.
func (m *mockDb) UpsertRecord(ctx context.Context, parent string, key string, value []byte) error {
	p, ok := m.records[parent]

	if !ok {
		m.records[parent] = make(map[string][]byte)
		p = m.records[parent]
	}

	p[key] = value

	return nil
}

func New() ports.Persistance {
	return &mockDb{
		records: make(map[string]map[string][]byte),
		audits:  make(map[string][][]byte),
	}
}

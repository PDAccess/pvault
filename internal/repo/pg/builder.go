package pg

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pdaccess/pvault/internal/core/domain"
	"github.com/pdaccess/pvault/internal/core/ports"
)

const (
	queryInsertRecord = `insert into pvault_kv (parent, key, value) values (:parent, :key, :value)`
	queryInsertAudit  = `insert into pvault_audit (key, audit) values (:key, :audit)`
	queryUpdateRecord = `update pvault_kv set value = :value where parent = :parent and key = :key`

	selectRecodrs = `select parent, key, value from pvault_kv`
	selectAudits  = `select key, audit from pvault_audit`
)

type PgPersistence struct {
	db *sqlx.DB
}

func New(connectionStr string) (ports.Persistance, error) {
	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("connection: %w", err)
	}

	if err := CreateSchema(context.Background(), connectionStr); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &PgPersistence{
		db: db,
	}, nil
}

// AppendAudit implements ports.Persistance.
func (p *PgPersistence) AppendAudit(ctx context.Context, key string, auditData []byte) error {
	_, err := p.db.NamedExecContext(ctx, queryInsertAudit, AuditRecord{
		KeyField:      key,
		OperationData: auditData,
	})

	if err != nil {
		return fmt.Errorf("AppendAudit: %w", err)
	}

	return nil
}

// UpsertRecord implements ports.Persistance.
func (p *PgPersistence) UpsertRecord(ctx context.Context, parent string, key string, value []byte) error {

	_, err := p.FetchRecord(ctx, parent, key)

	var query string
	if err == domain.ErrNotFound {
		query = queryInsertRecord
	} else {
		query = queryUpdateRecord
	}

	tx, err := p.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	_, err = p.db.NamedExecContext(ctx, query, KeyValueRecord{
		Parent:     parent,
		KeyField:   key,
		ValueField: value,
	})

	if err != nil {
		return fmt.Errorf("UpsertRecord: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// FetchRecord implements ports.Persistance.
func (p *PgPersistence) FetchRecord(ctx context.Context, parent string, key string) ([]byte, error) {
	query := selectRecodrs + " where parent=$1 and key=$2"
	stmt, err := p.db.PreparexContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("PrepareContext: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, parent, key)

	if err != nil {
		return nil, fmt.Errorf("QueryxContext: %w", err)
	}

	if !rows.Next() {
		return nil, domain.ErrNotFound
	}

	var kyre KeyValueRecord
	err = rows.StructScan(&kyre)
	if err != nil {
		return nil, fmt.Errorf("StructScan: %w", err)
	}

	return kyre.ValueField, nil
}

// SearchAudit implements ports.Persistance.
func (p *PgPersistence) SearchAudit(ctx context.Context, key string) ([][]byte, error) {
	query := selectAudits + " where key=$1"
	stmt, err := p.db.PreparexContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("PrepareContext: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, key)

	if err != nil {
		return nil, fmt.Errorf("QueryxContext: %w", err)
	}

	var ret [][]byte

	if rows.Next() {
		var audit AuditRecord
		err = rows.StructScan(&audit)
		if err != nil {
			return nil, fmt.Errorf("StructScan: %w", err)
		}

		ret = append(ret, audit.OperationData)

	}

	if len(ret) == 0 {
		return ret, domain.ErrNotFound
	}

	return ret, nil
}

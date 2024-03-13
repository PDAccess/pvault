package pg

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	createKvTable = `
		CREATE TABLE if not exists pvault_kv (
			parent text,
			key text,
			value bytea
		)`
	createAuditTable = `
		CREATE TABLE if not exists pvault_audit (
			key text,
			audit bytea
		)
		`
)

func CreateSchema(ctx context.Context, connectionStr string) error {
	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		return fmt.Errorf("sql connect: %w", err)
	}

	_, err = db.ExecContext(ctx, createKvTable)

	if err != nil {
		return fmt.Errorf("createKvTable: %w", err)
	}

	_, err = db.ExecContext(ctx, createAuditTable)

	if err != nil {
		return fmt.Errorf("createAuditTable: %w", err)
	}

	return nil
}

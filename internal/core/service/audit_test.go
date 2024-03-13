package service

import (
	"testing"

	"github.com/pdaccess/pvault/internal/core/domain"
)

func TestAduitLogs(t *testing.T) {

	key := domain.RecordKey("secrets/test")
	rec := domain.NewRecord()
	rec.KeyValueMap["password"] = "test123"

	err := impl.Write(ctx, key, rec)

	if err != nil {
		t.Fatalf("write shouldn't return an error: %v", err)
	}

	audits, err := impl.Logs(ctx, domain.NewRecordKey("test"))

	if err != nil {
		t.Fatalf("Write shouldn't return an error: %v", err)
	}

	if len(audits) == 0 {
		t.Fatalf("no audit logs")
	}

	if audits[0].OpType != domain.OpTypeWrite {
		t.Fatalf("wrong op type. given: %d want: %d", audits[0].OpType, domain.OpTypeWrite)
	}

	_, err = impl.Read(ctx, key)

	if err != nil {
		t.Fatalf("Read shouldn't return an error: %v", err)
	}

	audits, err = impl.Logs(ctx, domain.NewRecordKey("test"))

	if err != nil {
		t.Fatalf("Write shouldn't return an error: %v", err)
	}

	if len(audits) != 2 {
		t.Fatalf("no audit logs")
	}

	if audits[1].OpType != domain.OpTypeRead {
		t.Fatalf("wrong op type. given: %d want: %d", audits[0].OpType, domain.OpTypeRead)
	}
}

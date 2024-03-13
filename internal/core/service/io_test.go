package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pdaccess/pvault/internal/core/domain"
)

func TestWriteValidCases(t *testing.T) {
	key := domain.RecordKey("secrets/test")
	rec := domain.NewRecord()
	rec.KeyValueMap["password"] = "test123"

	err := impl.Write(ctx, key, rec)

	if err != nil {
		t.Errorf("write error: %v", err)
	}

	recRet, err := impl.Read(ctx, key)

	if err != nil {
		t.Errorf("read error: %v", err)
	}

	if err := compareRecords(rec, *recRet); err != nil {
		t.Errorf("return data is different. %v", err)
	}
}

func TestWriteUpsert(t *testing.T) {
	key := domain.RecordKey("secrets/test")
	rec := domain.NewRecord()
	rec.KeyValueMap["password"] = "test123"

	err := impl.Write(ctx, key, rec)

	if err != nil {
		t.Errorf("write error: %v", err)
	}

	rec.KeyValueMap["password"] = "test1234"

	err = impl.Write(ctx, key, rec)

	if err != nil {
		t.Errorf("write error: %v", err)
	}

	recRet, err := impl.Read(ctx, key)

	if err != nil {
		t.Errorf("read error: %v", err)
	}

	if err := compareRecords(rec, *recRet); err != nil {
		t.Errorf("return data is different. %v", err)
	}
}

func compareRecords(given, want domain.SimpleRecord) error {
	if !given.CreatedAt.Equal(want.CreatedAt) {
		return fmt.Errorf("wrong timestamps. given: %s, want: %s", given.CreatedAt, want.CreatedAt)
	}

	if !reflect.DeepEqual(given.KeyValueMap, want.KeyValueMap) {
		return fmt.Errorf("wrong kv map. given: %v, want: %v", given.KeyValueMap, want.KeyValueMap)
	}

	return nil
}

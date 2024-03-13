package pg

import (
	"context"
	"crypto/rand"
	"reflect"
	"testing"
)

func TestRecords(t *testing.T) {
	backend, err := New(connectionStr)

	if err != nil {
		t.Fatalf("backend shouldn't return error :%v", err)
	}

	buf := make([]byte, 128)
	_, err = rand.Read(buf)
	if err != nil {
		t.Fatalf("random shouldn't return error :%v", err)
	}

	err = backend.UpsertRecord(context.TODO(), "/sys", "/test", buf)

	if err != nil {
		t.Fatalf("UpsertRecord shouldn't return error :%v", err)
	}

	bytes, err := backend.FetchRecord(context.TODO(), "/sys", "/test")

	if err != nil {
		t.Fatalf("FetchRecord shouldn't return error :%v", err)
	}

	if !reflect.DeepEqual(buf, bytes) {
		t.Fatalf("Write and Read data are different")
	}
}

func TestAudits(t *testing.T) {
	backend, err := New(connectionStr)

	if err != nil {
		t.Errorf("backend shouldn't return error :%v", err)
	}

	buf := make([]byte, 128)
	_, err = rand.Read(buf)
	if err != nil {
		t.Fatalf("random shouldn't return error :%v", err)
	}

	err = backend.AppendAudit(context.TODO(), "/sys/test", buf)

	if err != nil {
		t.Fatalf("AppendAudit shouldn't return error :%v", err)
	}

	bytes, err := backend.SearchAudit(context.TODO(), "/sys/test")

	if err != nil {
		t.Fatalf("SearchAudit shouldn't return error :%v", err)
	}

	if !reflect.DeepEqual(buf, bytes[0]) {
		t.Fatalf("Write and Read data are different")
	}
}

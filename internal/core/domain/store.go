package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type KeyType string

const (
	KeyTypeUsername = KeyType("username")
	KeyTypePassword = KeyType("password")
	KeyTypeCert     = KeyType("cert")
	KeyTypeSshKey   = KeyType("ssh_key")
	KeyTypeOther    = KeyType("other")
)

type RecordKey string

func NewRecordKey(keyStr string) RecordKey {
	return RecordKey(keyStr)
}

func (r RecordKey) Split() (parent, children string) {
	strs := strings.Split(string(r), "/")
	return string(strs[0]), string(strs[1])
}

func (r RecordKey) String() string {
	p, k := r.Split()

	return fmt.Sprintf("parent: %s key: %s", p, k)
}

type SimpleRecord struct {
	Id          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	KeyValueMap map[string]string
}

func NewRecord() SimpleRecord {
	return SimpleRecord{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		KeyValueMap: make(map[string]string),
	}
}

func (s *SimpleRecord) Bytes() ([]byte, error) {
	return json.Marshal(*s)
}

func (s *SimpleRecord) Restore(raw []byte) error {
	return json.Unmarshal(raw, s)
}

func (s *SimpleRecord) Validate() error {
	return nil
}

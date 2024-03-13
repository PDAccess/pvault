package pg

type KeyValueRecord struct {
	Parent     string `db:"parent"`
	KeyField   string `db:"key"`
	ValueField []byte `db:"value"`
}

type AuditRecord struct {
	KeyField      string `db:"key"`
	OperationData []byte `db:"audit"`
}

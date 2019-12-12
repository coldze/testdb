package wraps

import "database/sql"

type Scanner interface {
	Closable
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
}

type ScannerFactory func(rows *sql.Rows) (Scanner, error)

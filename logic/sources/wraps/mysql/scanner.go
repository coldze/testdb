package mysql

import (
	"database/sql"
	"errors"

	"github.com/coldze/testdb/logic/sources/wraps"
)

type mysqlScanner struct {
	rows *sql.Rows
}

func (m *mysqlScanner) Scan(dest ...interface{}) error {
	return m.rows.Scan(dest...)
}

func (m *mysqlScanner) Next() bool {
	return m.rows.Next()
}

func (m *mysqlScanner) Err() error {
	return m.rows.Err()
}

func (m *mysqlScanner) Close() error {
	return m.rows.Close()
}

func NewScanner(rows *sql.Rows) (wraps.Scanner, error) {
	if rows == nil {
		return nil, errors.New("Rows is nil")
	}
	return &mysqlScanner{
		rows: rows,
	}, nil
}

package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/coldze/testdb/logic/sources/wraps"
)

type mysqlDBImpl struct {
	client     *sql.DB
	newScanner wraps.ScannerFactory
}

func (m *mysqlDBImpl) Query(q string, args ...interface{}) (wraps.Scanner, error) {
	rows, err := m.client.Query(q, args...)
	if err != nil {
		return nil, err
	}
	return m.newScanner(rows)
}

func (m *mysqlDBImpl) Close() error {
	return m.client.Close()
}

func NewMysqlDbWrap(connection string, factory wraps.ScannerFactory) (wraps.DbWrap, error) {
	client, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}
	return &mysqlDBImpl{
		client:     client,
		newScanner: factory,
	}, nil
}

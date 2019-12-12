package stores

import (
	"context"
	"errors"

	"github.com/coldze/testdb/logic/sources"
	"github.com/coldze/testdb/logic/sources/wraps"
	"github.com/coldze/testdb/logic/structs"
)

type dbDataSource struct {
	dbWrap   wraps.DbWrap
	qBuilder wraps.QueryBuilder
}

func (m *dbDataSource) Get(ctx context.Context, key structs.Request) (structs.Data, error) {
	q, err := m.qBuilder(key)
	if err != nil {
		return nil, err
	}
	scanner, err := m.dbWrap.Query(q.Request, q.Args...)
	if err != nil {
		return nil, err
	}
	if scanner == nil {
		return nil, errors.New("Result is nil")
	}
	var id string
	var count uint64
	res := structs.Data{}
	for scanner.Next() {
		err = scanner.Scan(&id, &count)
		if err != nil {
			return nil, err
		}
		res[id] = count
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewDbDataSource(wrap wraps.DbWrap, builder wraps.QueryBuilder) sources.Source {
	return &dbDataSource{
		dbWrap:   wrap,
		qBuilder: builder,
	}
}

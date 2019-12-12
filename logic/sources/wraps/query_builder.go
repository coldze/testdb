package wraps

import "github.com/coldze/testdb/logic/structs"

type Query struct {
	Request string
	Args    []interface{}
}

type QueryBuilder func(key structs.Request) (Query, error)

package sources

import (
	"context"
	"time"

	"github.com/coldze/testdb/logic/structs"
)

type Source interface {
	Get(ctx context.Context, key structs.Request) (structs.Data, error)
}

type Store interface {
	Put(ctx context.Context, date time.Time, data structs.Data) error
	Del(ctx context.Context, key structs.Request) error
	Wipe(ctx context.Context) error
}

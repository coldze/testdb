package caches

import (
	"net/http"

	"github.com/coldze/testdb/logic"
	"github.com/coldze/testdb/logic/handlers"
	"github.com/coldze/testdb/logic/sources"
)

func newWipeLogicHandler(cache sources.Store) handlers.LogicHandler {
	return func(r *http.Request) (interface{}, error) {
		return nil, cache.Wipe(r.Context())
	}
}

func NewWipeHandler(loggerFactory handlers.LoggerFactory, cache sources.Store) http.HandlerFunc {
	lHandler := newWipeLogicHandler(cache)
	wFactory := logic.DefaultWriterFactory()
	handler := handlers.NewHttpHandler(lHandler, wFactory)
	return handlers.NewCheckAndSetLoggerMiddleware(loggerFactory, handlers.GuidIDFactory, handler)
}

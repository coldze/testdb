package caches

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/coldze/testdb/logic"
	"github.com/coldze/testdb/logic/handlers"
	"github.com/coldze/testdb/logic/sources"
	"github.com/coldze/testdb/logic/structs"
)

func newDecodeQuery() (logic.Decoder, error) {
	reg, err := regexp.Compile(handlers.REGEX_PARSE_IDS)
	if err != nil {
		return nil, err
	}
	return func(r *http.Request) (res structs.Request, err error) {
		q := r.URL.Query()
		ids := q.Get(handlers.QUERY_KEY_ID)
		date := q.Get(handlers.QUERY_KEY_DATE)
		if len(date) <= 0 {
			err = errors.New("Date is required")
			return
		}
		parsedDate, err := time.Parse(handlers.EXPECTED_DATE_FORMAT, date)
		if err != nil {
			return
		}
		res.Date = parsedDate
		res.IDs = reg.FindAllString(ids, -1)
		if len(res.IDs) <= 0 {
			err = errors.New("ID is required")
		}
		return
	}, nil
}

func newDeleteLogicHandler(decode logic.Decoder, cache sources.Store) handlers.LogicHandler {
	return func(r *http.Request) (interface{}, error) {
		data, err := decode(r)
		if err != nil {
			return nil, err
		}
		return nil, cache.Del(r.Context(), data)
	}
}

func NewDeleteHandler(loggerFactory handlers.LoggerFactory, cache sources.Store) (http.HandlerFunc, error) {
	decoder, err := newDecodeQuery()
	if err != nil {
		return nil, err
	}
	lHandler := newDeleteLogicHandler(decoder, cache)
	wFactory := logic.DefaultWriterFactory()
	handler := handlers.NewHttpHandler(lHandler, wFactory)
	return handlers.NewCheckAndSetLoggerMiddleware(loggerFactory, handlers.GuidIDFactory, handler), nil
}

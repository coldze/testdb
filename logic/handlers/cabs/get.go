package cabs

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
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
		ignoreCache := q.Get(handlers.QUERY_KEY_NOCACHE)
		if len(ignoreCache) > 0 {
			res.IgnoreCache, err = strconv.ParseBool(ignoreCache)
			if err != nil {
				return
			}
		}
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

func newGetLogicHandler(decode logic.Decoder, source sources.Source) handlers.LogicHandler {
	return func(r *http.Request) (interface{}, error) {
		data, err := decode(r)
		if err != nil {
			return nil, err
		}
		return source.Get(r.Context(), data)
	}
}

func NewGetHandler(loggerFactory handlers.LoggerFactory, source sources.Source) (http.HandlerFunc, error) {
	decoder, err := newDecodeQuery()
	if err != nil {
		return nil, err
	}
	lHandler := newGetLogicHandler(decoder, source)
	wFactory := logic.DefaultWriterFactory()
	handler := handlers.NewHttpHandler(lHandler, wFactory)
	return handlers.NewCheckAndSetLoggerMiddleware(loggerFactory, handlers.GuidIDFactory, handler), nil
}

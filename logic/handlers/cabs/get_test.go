package cabs

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/mocks"
	"github.com/coldze/testdb/mocks/mock_handlers"
	"github.com/coldze/testdb/mocks/mock_logic"
	"github.com/coldze/testdb/mocks/mock_logs"
	"github.com/coldze/testdb/mocks/mock_sources"
	"github.com/coldze/testdb/utils"
	"github.com/golang/mock/gomock"
)

type getFixture struct {
	err error
	ctx context.Context
	r *http.Request
	res structs.Data
	logger    *mock_logs.MockLogger
	loggerFactory *mock_handlers.MockLoggerFactory
	source *mock_sources.MockSource
	decoder *mock_logic.MockDecoder
	key structs.Request
}

const (
	url_correct = "https://some.random.url/rest?id=123,345&date=2013-01-01"
	url_no_date = "https://some.random.url/rest?id=123,345"
	url_no_ids = "https://some.random.url/rest?date=2013-01-01"
	url_invalid_date = "https://some.random.url/rest?id=123,345&date=2013-33-01"
	url_invalid_ignore = "https://some.random.url/rest?id=123,345&date=2013-01-01&nocache=nope"
	url_correct_ignore = "https://some.random.url/rest?id=123,345&date=2013-01-01&nocache=true"
)

func newGetFixture(ctrl *gomock.Controller, url string, ignore bool) *getFixture {
	l := mock_logs.NewMockLogger(ctrl)
	ctx := utils.SetLogger(context.Background(), l)
	date, _ := time.Parse("2006-01-02", "2013-01-01")
	return &getFixture{
		loggerFactory: mock_handlers.NewMockLoggerFactory(ctrl),
		err: errors.New("Test"),
		ctx: ctx,
		r: httptest.NewRequest(http.MethodGet, url, nil),
		logger: l,
		source: mock_sources.NewMockSource(ctrl),
		decoder: mock_logic.NewMockDecoder(ctrl),
		res: structs.Data{
			"123": 99,
			"345": 6,
		},
		key: structs.Request{
			Date: date,
			IDs: []string{"123", "345"},
			IgnoreCache: ignore,
		},
	}
}

func TestCacheWipe(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_correct, false)

		f.decoder.EXPECT().Execute(f.r).Return(f.key, nil).Times(1)
		f.source.EXPECT().Get(f.r.Context(), f.key).Return(f.res, nil).Times(1)
		s := newGetLogicHandler(f.decoder.Execute, f.source)
		res, err := s(f.r)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(res, f.res) {
			t.FailNow()
		}

	})

	t.Run("fails if source fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_correct, false)

		f.decoder.EXPECT().Execute(f.r).Return(f.key, nil).Times(1)
		f.source.EXPECT().Get(f.r.Context(), f.key).Return(nil, f.err).Times(1)
		s := newGetLogicHandler(f.decoder.Execute, f.source)
		res, err := s(f.r)
		mocks.CmpError(t, err, f.err)
		resTyped, ok := res.(structs.Data)
		if !ok {
			t.FailNow()
		}
		if len(resTyped) > 0 {
			t.FailNow()
		}

	})

	t.Run("fails if decode fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_correct, false)

		f.decoder.EXPECT().Execute(f.r).Return(structs.Request{}, f.err).Times(1)
		s := newGetLogicHandler(f.decoder.Execute, f.source)
		res, err := s(f.r)
		mocks.CmpError(t, err, f.err)
		if res != nil {
			resTyped, ok := res.(structs.Data)
			if !ok {
				t.FailNow()
			}
			if len(resTyped) > 0 {
				t.FailNow()
			}
		}

	})
}

func TestDecode(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_correct, false)

		s, _ := newDecodeQuery()
		res, err := s(f.r)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(res, f.key) {
			t.FailNow()
		}

	})

	t.Run("success and ignore", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_correct_ignore, true)

		s, _ := newDecodeQuery()
		res, err := s(f.r)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(res, f.key) {
			t.FailNow()
		}

	})

	t.Run("fails if no ids", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_no_ids, false)

		s, _ := newDecodeQuery()
		_, err := s(f.r)
		if err == nil {
			t.FailNow()
		}

	})

	t.Run("fails if date fails to parse", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_invalid_date, false)

		s, _ := newDecodeQuery()
		_, err := s(f.r)
		if err == nil {
			t.FailNow()
		}

	})

	t.Run("fails if date is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_no_date, false)

		s, _ := newDecodeQuery()
		_, err := s(f.r)
		if err == nil {
			t.FailNow()
		}

	})

	t.Run("fails if nocache fails to parse", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newGetFixture(ctrl, url_invalid_ignore, false)

		s, _ := newDecodeQuery()
		_, err := s(f.r)
		if err == nil {
			t.FailNow()
		}

	})
}
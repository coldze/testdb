package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/utils"
	"github.com/coldze/testdb/mocks/mock_logs"
	"github.com/coldze/testdb/mocks/mock_std"
	"github.com/coldze/testdb/mocks/mock_handlers"
)

type middlewareFixture struct {
	Url       string
	RequestID string
	w         *mock_std.MockResponseWriter
	NewLogger *mock_handlers.MockLoggerFactory
	Logger    *mock_logs.MockLogger
	Next      *mock_std.MockHttpHandlerFunc
	Reader    *mock_std.MockReadCloser
	NewID     *mock_handlers.MockIDFactory
}

func newMiddlewareFixture(ctrl *gomock.Controller) *middlewareFixture {
	return &middlewareFixture{
		Url:       "https://test.url.com/",
		RequestID: "123-test-id",
		w:         mock_std.NewMockResponseWriter(ctrl),
		NewLogger: mock_handlers.NewMockLoggerFactory(ctrl),
		Logger:    mock_logs.NewMockLogger(ctrl),
		Next:      mock_std.NewMockHttpHandlerFunc(ctrl),
		Reader:    mock_std.NewMockReadCloser(ctrl),
		NewID:     mock_handlers.NewMockIDFactory(ctrl),
	}
}

func TestHttpMiddleware(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newMiddlewareFixture(ctrl)

		next := func(w http.ResponseWriter, r *http.Request) {
			f.Next.Handle(w, r)
			logger := utils.GetLogger(r.Context())
			if logger != f.Logger {
				t.Errorf("Not expected logger value: %v", logger)
			}
		}

		handle := NewCheckAndSetLoggerMiddleware(f.NewLogger.Create, f.NewID.Create, next)
		r := httptest.NewRequest(http.MethodGet, "https://test.url.com/", f.Reader)

		f.NewID.EXPECT().Create().Return(f.RequestID).Times(1)
		f.NewLogger.EXPECT().Create(f.RequestID).Return(f.Logger).Times(1)
		f.Next.EXPECT().Create(f.w, gomock.Any()).Times(1)

		handle(f.w, r)
	})

	t.Run("failed if request is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newMiddlewareFixture(ctrl)

		next := func(w http.ResponseWriter, r *http.Request) {
			f.Next.Handle(w, r)
			logger := utils.GetLogger(r.Context())
			if logger != f.Logger {
				t.Errorf("Not expected logger value: %v", logger)
			}
		}

		handle := NewCheckAndSetLoggerMiddleware(f.NewLogger.Create, f.NewID.Create, next)

		f.w.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

		handle(f.w, nil)
	})
}

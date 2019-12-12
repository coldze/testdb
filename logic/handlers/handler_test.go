package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/logic"
	"github.com/coldze/testdb/mocks/mock_handlers"
	"github.com/coldze/testdb/mocks/mock_logs"
	"github.com/coldze/testdb/mocks/mock_std"
	"github.com/coldze/testdb/utils"
)

type handlerFixture struct {
	RequestID     string
	Url           string
	Ctx           context.Context
	CtxWithHeader context.Context
	Data          string
	Error         error
	Logger        *mock_logs.MockLogger
	LogicHandler  *mock_handlers.MockLogicHandler
	Writer        *mock_handlers.MockWriter
	W             *mock_std.MockResponseWriter
	Response      interface{}
}

func newHandlerFixture(ctrl *gomock.Controller) *handlerFixture {
	reqID := "123-test-request-id"
	logger := mock_logs.NewMockLogger(ctrl)
	ctx := utils.SetLogger(context.Background(), logger)
	ctxWithHeader := utils.SetRequestID(ctx, reqID)

	return &handlerFixture{
		RequestID:     reqID,
		Url:           "https://test.com.au/",
		Ctx:           ctx,
		CtxWithHeader: ctxWithHeader,
		Data:          "some random data",
		Error:         errors.New("Some test error"),
		Logger:        logger,
		LogicHandler:  mock_handlers.NewMockLogicHandler(ctrl),
		Writer:        mock_handlers.NewMockWriter(ctrl),
		W:             mock_std.NewMockResponseWriter(ctrl),
	}
}

func newRequest(ctx context.Context, method string, url string, body io.Reader, headers http.Header) *http.Request {
	r := httptest.NewRequest(method, url, body)
	r.Header = headers
	return r.WithContext(ctx)
}

func newTestableHttpHandler(f *handlerFixture) http.HandlerFunc {
	return NewHttpHandler(f.LogicHandler.Handle, func(ctx context.Context, w http.ResponseWriter) logic.Writer {
		return f.Writer
	})
}

func TestHttpHandler(t *testing.T) {
	t.Run("error while getting data is a failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)
		r := httptest.NewRequest(http.MethodGet, f.Url, nil)
		r = r.WithContext(f.Ctx)

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(r).Return(nil, f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.Writer.EXPECT().Error(gomock.Any())

		httpHandler(f.W, r)
	})

	t.Run("error while getting data is a failure, failure writing error to response is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)
		r := httptest.NewRequest(http.MethodGet, f.Url, nil)
		r = r.WithContext(f.Ctx)

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(r).Return(nil, f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.Writer.EXPECT().Error(gomock.Any()).Return(f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("error during writing is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)
		r := httptest.NewRequest(http.MethodGet, f.Url, nil)
		r = r.WithContext(f.Ctx)

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(r).Return(make(chan int), nil).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		f.Writer.EXPECT().Data(gomock.Any()).Return(f.Error).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("writes error on throw of error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(gomock.Any()).Do(func(a interface{}) (interface{}, error) {
			panic(f.Error)
			return nil, nil
		}).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.Writer.EXPECT().Error(f.Error).Return(nil).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("writes error on throw of error, logs error if failed to write a response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(gomock.Any()).Do(func(a interface{}) (interface{}, error) {
			panic(f.Error)
			return nil, nil
		}).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)
		f.Writer.EXPECT().Error(f.Error).Return(f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("writes error on throw of unknown", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		unknownErr := "unknown error type"
		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(gomock.Any()).Do(func(a interface{}) (interface{}, error) {
			panic(unknownErr)
			return nil, nil
		}).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), unknownErr, unknownErr)
		f.Writer.EXPECT().Error(gomock.Any()).Return(nil).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("writes error on throw of unknown, logs error if failed to write a response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		unknownErr := "unknown error type"
		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, nil, http.Header{})

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(gomock.Any()).Do(func(a interface{}) (interface{}, error) {
			panic(unknownErr)
			return nil, nil
		}).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), unknownErr, unknownErr)
		f.Writer.EXPECT().Error(gomock.Any()).Return(f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("if body not nil, it gets closed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)
		body := mock_std.NewMockReadCloser(ctrl)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, body, http.Header{})

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(gomock.Any()).Return(nil, nil).Times(1)
		f.Writer.EXPECT().Data(nil).Return(nil).Times(1)

		body.EXPECT().Close().Return(nil).Times(1)

		httpHandler(f.W, r)
	})

	t.Run("if body returns error when closed, it gets logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newHandlerFixture(ctrl)
		httpHandler := newTestableHttpHandler(f)
		body := mock_std.NewMockReadCloser(ctrl)

		r := newRequest(f.Ctx, http.MethodGet, f.Url, body, http.Header{})

		f.Logger.EXPECT().Infof("Request URL: %v", f.Url)
		f.LogicHandler.EXPECT().Handle(gomock.Any()).Return(nil, nil).Times(1)
		f.Writer.EXPECT().Data(nil).Return(nil).Times(1)

		body.EXPECT().Close().Return(f.Error).Times(1)
		f.Logger.EXPECT().Errorf(gomock.Any(), f.Error).Times(1)

		httpHandler(f.W, r)
	})
}

func TestNewHttpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := newHandlerFixture(ctrl)
	httpHandler := NewHttpHandler(f.LogicHandler.Handle, func(ctx context.Context, w http.ResponseWriter) logic.Writer {
		return f.Writer
	})
	if httpHandler == nil {
		t.Errorf("Factory returns nil")
	}
}

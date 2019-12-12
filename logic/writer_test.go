package logic

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/mocks"
	"github.com/coldze/testdb/mocks/mock_logic"
	"github.com/coldze/testdb/mocks/mock_logs"
	"github.com/coldze/testdb/mocks/mock_std"
	"github.com/coldze/testdb/utils"
)

type writerFixture struct {
	ctx     context.Context
	encoder *mock_logic.MockEncoder
	w       *mock_std.MockResponseWriter
	e       error
	l       *mock_logs.MockLogger
}

func newWriterFixture(ctrl *gomock.Controller) *writerFixture {
	l := mock_logs.NewMockLogger(ctrl)
	ctx := utils.SetLogger(context.Background(), l)
	return &writerFixture{
		w:       mock_std.NewMockResponseWriter(ctrl),
		l:       l,
		ctx:     ctx,
		encoder: mock_logic.NewMockEncoder(ctrl),
		e:       errors.New("123"),
	}
}

func TestWriterFactory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		w := DefaultWriterFactory()
		respWriter := mock_std.NewMockResponseWriter(ctrl)
		w(context.Background(), respWriter)
	})
}

func TestWriterImpl(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newWriterFixture(ctrl)
		w := writerImpl{
			ctx:        f.ctx,
			marshal:    f.encoder.Execute,
			respWriter: f.w,
		}

		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, nil).Times(1)
		f.w.EXPECT().WriteHeader(http.StatusOK).Times(1)
		f.w.EXPECT().Header().Return(http.Header{}).Times(1)
		f.w.EXPECT().Write(gomock.Any()).Return(10, nil).Times(1)

		err := w.Data("123")
		if err != nil {
			t.FailNow()
		}
	})

	t.Run("returns error if write data failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newWriterFixture(ctrl)
		w := writerImpl{
			ctx:        f.ctx,
			marshal:    f.encoder.Execute,
			respWriter: f.w,
		}

		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, nil).Times(1)
		f.w.EXPECT().WriteHeader(http.StatusOK).Times(1)
		f.w.EXPECT().Header().Return(http.Header{}).Times(1)
		f.w.EXPECT().Write(gomock.Any()).Return(0, f.e).Times(1)

		err := w.Data("123")
		mocks.CmpError(t, err, f.e)
	})

	t.Run("if error log and write error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newWriterFixture(ctrl)
		w := writerImpl{
			ctx:        f.ctx,
			marshal:    f.encoder.Execute,
			respWriter: f.w,
		}

		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, f.e).Times(1)
		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, nil).Times(1)
		f.l.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		f.w.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)
		f.w.EXPECT().Header().Return(http.Header{}).Times(1)
		f.w.EXPECT().Write(gomock.Any()).Return(0, nil).Times(1)

		err := w.Data("123")
		if err != nil {
			t.FailNow()
		}
	})

	t.Run("if error log and write error, if write error fails, returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newWriterFixture(ctrl)
		w := writerImpl{
			ctx:        f.ctx,
			marshal:    f.encoder.Execute,
			respWriter: f.w,
		}

		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, f.e).Times(1)
		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, nil).Times(1)
		f.l.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		f.w.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)
		f.w.EXPECT().Header().Return(http.Header{}).Times(1)
		f.w.EXPECT().Write(gomock.Any()).Return(0, f.e).Times(1)

		err := w.Data("123")
		mocks.CmpError(t, err, f.e)
	})

	t.Run("write error success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newWriterFixture(ctrl)
		w := writerImpl{
			ctx:        f.ctx,
			marshal:    f.encoder.Execute,
			respWriter: f.w,
		}

		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, nil).Times(1)
		f.w.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)
		f.w.EXPECT().Header().Return(http.Header{}).Times(1)
		f.w.EXPECT().Write(gomock.Any()).Return(0, nil).Times(1)

		err := w.Error(errors.New("test"))
		if err != nil {
			t.FailNow()
		}
	})

	t.Run("write error fails, if encoding fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newWriterFixture(ctrl)
		w := writerImpl{
			ctx:        f.ctx,
			marshal:    f.encoder.Execute,
			respWriter: f.w,
		}

		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, f.e).Times(1)
		f.w.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

		err := w.Error(errors.New("test"))
		mocks.CmpError(t, err, f.e)
	})

	t.Run("write error fails, if write fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newWriterFixture(ctrl)
		w := writerImpl{
			ctx:        f.ctx,
			marshal:    f.encoder.Execute,
			respWriter: f.w,
		}

		f.encoder.EXPECT().Execute(f.ctx, gomock.Any()).Return([]byte{}, nil).Times(1)
		f.w.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)
		f.w.EXPECT().Header().Return(http.Header{}).Times(1)
		f.w.EXPECT().Write(gomock.Any()).Return(0, f.e).Times(1)

		err := w.Error(errors.New("test"))
		mocks.CmpError(t, err, f.e)
	})
}

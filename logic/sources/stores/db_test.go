package stores

import (
	"context"
	"errors"
	"reflect"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/logic/sources/wraps"
	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/mocks"
	"github.com/coldze/testdb/mocks/mock_logs"
	"github.com/coldze/testdb/mocks/mock_wraps"
	"github.com/coldze/testdb/utils"
)

type dbFixture struct {
	err error
	key structs.Request
	ctx context.Context
	logger    *mock_logs.MockLogger
	dbWrap *mock_wraps.MockDbWrap
	qBuilder *mock_wraps.MockQueryBuilder
	scanner *mock_wraps.MockScanner
}

func newDbFixture(ctrl *gomock.Controller) *dbFixture {
	l := mock_logs.NewMockLogger(ctrl)
	ctx := utils.SetLogger(context.Background(), l)
	return &dbFixture{
		err: errors.New("Test"),
		key: structs.Request{
			IDs: []string{"123", "2345", "123"},
			IgnoreCache: true,
			Date: time.Now(),
		},
		ctx: ctx,
		logger: l,
		dbWrap: mock_wraps.NewMockDbWrap(ctrl),
		qBuilder: mock_wraps.NewMockQueryBuilder(ctrl),
		scanner: mock_wraps.NewMockScanner(ctrl),

	}
}

func TestDbDataSource(t *testing.T) {

	t.Run("fails if query builder fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newDbFixture(ctrl)
		s := NewDbDataSource(f.dbWrap, f.qBuilder.Build)

		f.qBuilder.EXPECT().Execute(f.key).Return(wraps.Query{}, f.err).Times(1)

		_, err := s.Get(f.ctx, f.key)

		mocks.CmpError(t, err, f.err)

	})

	t.Run("fails if original fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newDbFixture(ctrl)
		s := NewDbDataSource(f.dbWrap, f.qBuilder.Build)

		f.qBuilder.EXPECT().Execute(f.key).Return(wraps.Query{}, nil).Times(1)
		f.dbWrap.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, f.err).Times(1)

		_, err := s.Get(f.ctx, f.key)

		mocks.CmpError(t, err, f.err)

	})

	t.Run("fails if scanner is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newDbFixture(ctrl)
		s := NewDbDataSource(f.dbWrap, f.qBuilder.Build)

		f.qBuilder.EXPECT().Execute(f.key).Return(wraps.Query{}, nil).Times(1)
		f.dbWrap.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

		_, err := s.Get(f.ctx, f.key)

		if err == nil {
			t.FailNow()
		}
	})

	t.Run("fails if scanner has error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newDbFixture(ctrl)
		s := NewDbDataSource(f.dbWrap, f.qBuilder.Build)

		f.qBuilder.EXPECT().Execute(f.key).Return(wraps.Query{}, nil).Times(1)
		f.dbWrap.EXPECT().Query(gomock.Any(), gomock.Any()).Return(f.scanner, nil).Times(1)
		f.scanner.EXPECT().Next().Return(false).Times(1)
		f.scanner.EXPECT().Err().Return(f.err).Times(1)

		_, err := s.Get(f.ctx, f.key)

		mocks.CmpError(t, err, f.err)

	})

	t.Run("ok if scanner has no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newDbFixture(ctrl)
		s := NewDbDataSource(f.dbWrap, f.qBuilder.Build)

		f.qBuilder.EXPECT().Execute(f.key).Return(wraps.Query{}, nil).Times(1)
		f.dbWrap.EXPECT().Query(gomock.Any(), gomock.Any()).Return(f.scanner, nil).Times(1)
		f.scanner.EXPECT().Next().Return(false).Times(1)
		f.scanner.EXPECT().Err().Return(nil).Times(1)

		_, err := s.Get(f.ctx, f.key)

		if err != nil {
			t.FailNow()
		}

	})

	t.Run("runs until scanner is done", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newDbFixture(ctrl)
		s := NewDbDataSource(f.dbWrap, f.qBuilder.Build)

		f.qBuilder.EXPECT().Execute(f.key).Return(wraps.Query{}, nil).Times(1)
		f.dbWrap.EXPECT().Query(gomock.Any(), gomock.Any()).Return(f.scanner, nil).Times(1)
		f.scanner.EXPECT().Next().Return(true).Times(4)
		f.scanner.EXPECT().Next().Return(false).Times(1)
		iter := uint64(1)
		res := structs.Data{}
		f.scanner.EXPECT().Scan(gomock.Any(), gomock.Any()).Do(func(a *string, b *uint64) {
			*a = fmt.Sprintf("id_%v", iter)
			*b = uint64(iter)
			res[*a] = *b
			iter++
		}).Return(nil).Times(4)
		f.scanner.EXPECT().Err().Return(nil).Times(1)

		r, err := s.Get(f.ctx, f.key)

		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(r, res) {
			t.FailNow()
		}

	})

	t.Run("fails if scanner fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newDbFixture(ctrl)
		s := NewDbDataSource(f.dbWrap, f.qBuilder.Build)

		f.qBuilder.EXPECT().Execute(f.key).Return(wraps.Query{}, nil).Times(1)
		f.dbWrap.EXPECT().Query(gomock.Any(), gomock.Any()).Return(f.scanner, nil).Times(1)
		f.scanner.EXPECT().Next().Return(true).Times(3)
		iter := uint64(1)
		res := structs.Data{}
		f.scanner.EXPECT().Scan(gomock.Any(), gomock.Any()).Do(func(a *string, b *uint64) {
			*a = fmt.Sprintf("id_%v", iter)
			*b = uint64(iter)
			res[*a] = *b
			iter++
		}).Return(nil).Times(2)
		f.scanner.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(f.err).Times(1)

		_, err := s.Get(f.ctx, f.key)

		mocks.CmpError(t, err, f.err)

	})
}


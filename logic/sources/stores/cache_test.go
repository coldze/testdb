package stores

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/go-redis/redis"

	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/mocks"
	"github.com/coldze/testdb/mocks/mock_logs"
	"github.com/coldze/testdb/mocks/mock_wraps"
	"github.com/coldze/testdb/utils"
)

type cacheFixture struct {
	err error
	key structs.Request
	data structs.Data
	ctx context.Context
	logger    *mock_logs.MockLogger
	cacheWrap *mock_wraps.MockCacheWrap
	ttl time.Duration
}

func newCacheFixture(ctrl *gomock.Controller) *cacheFixture {
	l := mock_logs.NewMockLogger(ctrl)
	ctx := utils.SetLogger(context.Background(), l)
	return &cacheFixture{
		err: errors.New("Test"),
		key: structs.Request{
			IDs: []string{"123", "2345", "1238"},
			IgnoreCache: true,
			Date: time.Now(),
		},
		data: structs.Data{
			"123": 5,
			"2345": 35,
			"1238": 1,
		},
		ctx: ctx,
		logger: l,
		cacheWrap: mock_wraps.NewMockCacheWrap(ctrl),
		ttl: time.Second,
	}
}

func TestCacheWipe(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		f.cacheWrap.EXPECT().Flush().Return(nil).Times(1)
		err := s.Wipe(f.ctx)
		if err != nil {
			t.FailNow()
		}

	})

	t.Run("fails, if wrap fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		f.cacheWrap.EXPECT().Flush().Return(f.err).Times(1)
		err := s.Wipe(f.ctx)
		mocks.CmpError(t, err, f.err)
	})
}

func TestCachePut(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, v := range f.data {
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Set(key, v, f.ttl).Return(nil).Times(1)
		}

		err := s.Put(f.ctx, f.key.Date, f.data)
		if err != nil {
			t.FailNow()
		}

	})

	t.Run("logs warning if cache fails to set	", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, v := range f.data {
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Set(key, v, f.ttl).Return(f.err).Times(1)
			f.logger.EXPECT().Warningf(gomock.Any(), key, f.err).Times(1)
		}

		err := s.Put(f.ctx, f.key.Date, f.data)
		if err != nil {
			t.FailNow()
		}

	})

}

func TestCacheDel(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for _, k := range f.key.IDs {
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Del(key).Return(nil).Times(1)
		}

		err := s.Del(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}

	})

	t.Run("logs warning if cache fails to set	", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for _, k := range f.key.IDs {
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Del(key).Return(f.err).Times(1)
			f.logger.EXPECT().Warningf(gomock.Any(), key, f.err).Times(1)
		}

		err := s.Del(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}

	})

}

func TestCacheGet(t *testing.T) {

	t.Run("if key has no ids, result is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		f.key.IDs = []string{}
		s := NewCache(f.cacheWrap, f.ttl)

		res, err := s.Get(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}
		if len(res) > 0 {
			t.FailNow()
		}

	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, v := range f.data{
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Get(key).Return(fmt.Sprintf("%v", v), nil).Times(1)
		}

		res, err := s.Get(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(res, f.data) {
			t.FailNow()
		}

	})

	t.Run("value skipped if cache returns Nil-error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, _ := range f.data{
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Get(key).Return("", redis.Nil).Times(1)
		}

		res, err := s.Get(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}
		if len(res) > 0 {
			t.FailNow()
		}

	})

	t.Run("value skipped if cache returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, _ := range f.data{
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Get(key).Return(nil, nil).Times(1)
		}

		res, err := s.Get(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}
		if len(res) > 0 {
			t.FailNow()
		}

	})

	t.Run("warning logged if cache returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, _ := range f.data{
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Get(key).Return(nil, f.err).Times(1)
			f.logger.EXPECT().Warningf(gomock.Any(), key, f.err)
		}

		res, err := s.Get(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}
		if len(res) > 0 {
			t.FailNow()
		}

	})

	t.Run("warning logged if cache returns not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, _ := range f.data{
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Get(key).Return(123, nil).Times(1)
			f.logger.EXPECT().Warningf(gomock.Any(), key, 123)
		}

		res, err := s.Get(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}
		if len(res) > 0 {
			t.FailNow()
		}

	})

	t.Run("warning logged if cache returns a string that is not convertable to uint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCacheFixture(ctrl)
		s := NewCache(f.cacheWrap, f.ttl)

		for k, _ := range f.data{
			key := buildCacheKey(k, f.key.Date)
			f.cacheWrap.EXPECT().Get(key).Return("-123", nil).Times(1)
			f.logger.EXPECT().Warningf(gomock.Any(), key, gomock.Any())
		}

		res, err := s.Get(f.ctx, f.key)
		if err != nil {
			t.FailNow()
		}
		if len(res) > 0 {
			t.FailNow()
		}

	})
}

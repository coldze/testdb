package sources

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/mocks"
	"github.com/coldze/testdb/mocks/mock_logs"
	"github.com/coldze/testdb/mocks/mock_sources"
	"github.com/coldze/testdb/utils"
)

type cachedFixture struct {
	cache         *mock_sources.MockCache
	source *mock_sources.MockSource
	logger    *mock_logs.MockLogger
	ctx context.Context
	req structs.Request
	res structs.Data
	err error
}

func newCachedFixture(ctrl *gomock.Controller, ignore bool) *cachedFixture {
	l := mock_logs.NewMockLogger(ctrl)
	ctx := utils.SetLogger(context.Background(), l)
	return &cachedFixture{
		cache: mock_sources.NewMockCache(ctrl),
		source: mock_sources.NewMockSource(ctrl),
		logger: l,
		ctx: ctx,
		err: errors.New("ssdd"),
		res: structs.Data{
			"123": 8,
			"2134": 18,
		},
		req: structs.Request{
			IDs: []string{"123", "2134"},
			IgnoreCache: ignore,
			Date: time.Now(),
		},
	}
}

func TestCachedIgnoreCache(t *testing.T) {

	t.Run("fails if original fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCachedFixture(ctrl, true)
		c := NewSourceWithCache(f.source, f.cache)

		f.source.EXPECT().Get(f.ctx, f.req).Return(nil, f.err).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), structs.Data{}).Times(1)

		_, err := c.Get(f.ctx, f.req)

		mocks.CmpError(t, err, f.err)

	})

	t.Run("logs error if fails to put to cache", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCachedFixture(ctrl, true)
		c := NewSourceWithCache(f.source, f.cache)

		f.source.EXPECT().Get(f.ctx, f.req).Return(f.res, nil).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)
		f.cache.EXPECT().Put(f.ctx, f.req.Date, f.res).Return(f.err).Times(1)
		f.logger.EXPECT().Warningf(gomock.Any(), f.err).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), f.res).Times(1)


		res, err := c.Get(f.ctx, f.req)
		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(res, f.res) {
			t.FailNow()
		}
	})
}

func TestCachedUseCache(t *testing.T) {

	t.Run("logs error if cache fails, falls through to original and fails if failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCachedFixture(ctrl, false)
		c := NewSourceWithCache(f.source, f.cache)

		f.cache.EXPECT().Get(f.ctx, f.req).Return(nil, f.err)
		f.logger.EXPECT().Warningf(gomock.Any(), f.err).Times(1)

		f.source.EXPECT().Get(f.ctx, f.req).Return(nil, f.err).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), structs.Data{}).Times(1)

		_, err := c.Get(f.ctx, f.req)

		mocks.CmpError(t, err, f.err)

	})

	t.Run("if nothing in cache, falls through to original and fails if failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCachedFixture(ctrl, false)
		c := NewSourceWithCache(f.source, f.cache)

		f.cache.EXPECT().Get(f.ctx, f.req).Return(nil, nil)

		f.source.EXPECT().Get(f.ctx, f.req).Return(nil, f.err).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), structs.Data{}).Times(1)

		_, err := c.Get(f.ctx, f.req)

		mocks.CmpError(t, err, f.err)

	})

	t.Run("if part is in cache, falls through to original fills rest with data from original", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCachedFixture(ctrl, false)
		c := NewSourceWithCache(f.source, f.cache)
		dataInCache := structs.Data{
			f.req.IDs[0]: 999,
		}
		dataInSource := structs.Data{
			f.req.IDs[1]: 234,
		}
		f.res[f.req.IDs[0]] = 999
		f.res[f.req.IDs[1]] = 234
		srcReq := structs.Request{
			IDs: []string{f.req.IDs[1]},
			Date: f.req.Date,
			IgnoreCache: f.req.IgnoreCache,
		}

		f.cache.EXPECT().Get(f.ctx, f.req).Return(dataInCache, nil)
		f.logger.EXPECT().Infof(gomock.Any(), dataInCache).Times(1)

		f.source.EXPECT().Get(f.ctx, srcReq).Return(dataInSource, nil).Times(1)
		f.cache.EXPECT().Put(f.ctx, f.req.Date, dataInSource).Return(nil).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), dataInSource).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)




		res, err := c.Get(f.ctx, f.req)

		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(res, f.res) {
			t.FailNow()
		}

	})

	t.Run("if everything is in cache, returns cache data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := newCachedFixture(ctrl, false)
		c := NewSourceWithCache(f.source, f.cache)

		f.cache.EXPECT().Get(f.ctx, f.req).Return(f.res, nil)
		f.logger.EXPECT().Infof(gomock.Any(), f.res).Times(1)
		f.logger.EXPECT().Infof(gomock.Any(), f.res).Times(1)


		res, err := c.Get(f.ctx, f.req)

		if err != nil {
			t.FailNow()
		}
		if !reflect.DeepEqual(res, f.res) {
			t.FailNow()
		}

	})
}

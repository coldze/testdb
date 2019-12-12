package logic

import (
	"context"
	"math"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/logic/structs"
)

func TestEncodeResponse(t *testing.T) {
	t.Run("error is returned if fails to marshal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		_, err := EncodeResponse(ctx, structs.ResponseData{
			Data: math.Inf(-1),
		})
		if err == nil {
			t.FailNow()
		}
	})

	t.Run("ok if marshal is ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		_, err := EncodeResponse(ctx, structs.ResponseData{
			Data: "123",
		})
		if err != nil {
			t.FailNow()
		}
	})
}

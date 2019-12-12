package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestIDFactory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		GuidIDFactory()
	})
}

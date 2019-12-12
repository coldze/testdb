package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/coldze/testdb/mocks/mock_logs"
)

func TestDefaultLoggerFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logs.NewMockLogger(ctrl)
	lf := NewDefaultLoggerFactory(logger)
	l := lf("123")
	logger.EXPECT().Infof(" [123] info-test").Times(1)
	logger.EXPECT().Warningf(" [123] warning-test").Times(1)
	logger.EXPECT().Debugf(" [123] debug-test").Times(1)
	logger.EXPECT().Errorf(" [123] error-test").Times(1)

	l.Infof("info-test")
	l.Warningf("warning-test")
	l.Debugf("debug-test")
	l.Errorf("error-test")

}

package logic

import (
	"context"
	"net/http"

	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/utils"
)

const (
	HEADER_CONTENT_TYPE   = "Content-Type"
	MIME_APPLICATION_JSON = "application/json"
)

type Writer interface {
	Error(err error) error
	Data(res interface{}) error
}

type writerImpl struct {
	ctx        context.Context
	respWriter http.ResponseWriter
	marshal    Encoder
}

func (w *writerImpl) setType() {
	w.respWriter.Header().Set(HEADER_CONTENT_TYPE, MIME_APPLICATION_JSON)
}

func (w *writerImpl) Error(res error) error {
	w.respWriter.WriteHeader(http.StatusInternalServerError)
	data, err := w.marshal(w.ctx, structs.ResponseData{
		Error: res.Error(),
	})
	if err != nil {
		return err
	}
	w.setType()
	_, err = w.respWriter.Write(data)
	return err
}

func (w *writerImpl) Data(res interface{}) error {
	data, err := w.marshal(w.ctx, structs.ResponseData{
		Data: res,
	})
	if err != nil {
		logger := utils.GetLogger(w.ctx)
		logger.Errorf("Failed to write response. Error: %v", err)
		return w.Error(err)
	}
	w.respWriter.WriteHeader(http.StatusOK)
	w.setType()
	_, err = w.respWriter.Write(data)
	return err
}

type WriterFactory func(ctx context.Context, w http.ResponseWriter) Writer

func DefaultWriterFactory() WriterFactory {
	return func(ctx context.Context, w http.ResponseWriter) Writer {
		return &writerImpl{
			marshal:    EncodeResponse,
			ctx:        ctx,
			respWriter: w,
		}
	}
}

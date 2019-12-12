package handlers

import (
	"fmt"
	"net/http"

	"github.com/coldze/testdb/utils"
	"github.com/coldze/testdb/logic"
)

type LogicHandler func(r *http.Request) (interface{}, error)

func NewHttpHandler(handler LogicHandler, newWriter logic.WriterFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := utils.GetLogger(ctx)
		logger.Infof("Request URL: %v", r.URL.String())
		writer := newWriter(ctx, w)
		body := r.Body
		if body != nil {
			defer func() {
				err := body.Close()
				if err != nil {
					logger.Errorf("Failed to close body. Error: %v", err)
				}
			}()
		}
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			err, ok := r.(error)
			if ok {
				logger.Errorf("Panic occurred in handler. Error: %v", err)
			} else {
				logger.Errorf("Panic occurred in handler. Unknown error: %+v. Type: %T.", r, r)
				err = fmt.Errorf("Panic occurred in handler. Unknown error: %+v. Type: %T.", r, r)
			}
			wErr := writer.Error(err)
			if wErr != nil {
				logger.Errorf("Failed to write error to response. Error: %v", wErr)
			}
		}()
		data, err := handler(r)
		if err != nil {
			logger.Errorf("Failed to read body. Error: %v", err)
			wErr := writer.Error(err)
			if wErr != nil {
				logger.Errorf("Failed to write error to response. Error: %v", wErr)
			}
			return
		}
		err = writer.Data(data)
		if err != nil {
			//unexpected behaviour, failed to write response - seems like trying to write error is no use then. Log error and quit func.
			logger.Errorf("Failed to write response. Error: %v", err)
		}
	}
}

package logic

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/utils"
)

type Decoder func(r *http.Request) (structs.Request, error)
type Encoder func(ctx context.Context, res structs.ResponseData) ([]byte, error)

func EncodeResponse(ctx context.Context, res structs.ResponseData) ([]byte, error) {
	rid := utils.GetRequestID(ctx)
	r := structs.Response{
		ResponseData: res,
		RequestID:    rid,
	}
	return json.Marshal(r)
}

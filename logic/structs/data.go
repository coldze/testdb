package structs

import "time"

type Request struct {
	IDs         []string
	Date        time.Time
	IgnoreCache bool
}

type Data map[string]uint64

type ResponseData struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data"`
}

type Response struct {
	RequestID string `json:"trace_id"`
	ResponseData
}

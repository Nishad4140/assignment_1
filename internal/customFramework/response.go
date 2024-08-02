package customframework

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	headerWritten bool
}

func (r *Response) Status(HTTPStatus int) *Response {
	if !r.headerWritten {
		r.WriteHeader(HTTPStatus)
		r.headerWritten = true
	}
	return r
}

func (r *Response) Header(header, value string) *Response {
	r.ResponseWriter.Header().Set(header, value)
	return r
}

func (r *Response) End() {
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
	}
}

func (r *Response) Json(value any) {
	r.Header("Content-Type", "application/json")
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
		r.headerWritten = true
	}
	json.NewEncoder(r.ResponseWriter).Encode(value)
}

func (r *Response) Write(data []byte) *Response {
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
		r.headerWritten = true
	}
	r.ResponseWriter.Write(data)
	return r
}

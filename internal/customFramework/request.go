package customframework

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	*http.Request
	params map[string]string
}

func (r *Request) Headers() map[string]any {
	headers := make(map[string]any)
	for key, values := range r.Header {
		if len(values) == 1 {
			headers[key] = values[0]
		} else {
			headers[key] = values
		}
	}
	return headers
}

func (r *Request) Query() map[string]any {
	queries := make(map[string]any)
	for key, values := range r.URL.Query() {
		if len(values) == 1 {
			queries[key] = values[0]
		} else {
			queries[key] = values
		}
	}
	return queries
}

func (r *Request) PathParam(param string) string {
	return r.params[param]
}

func (r *Request) Body(v any) error {
	return json.NewDecoder(r.Request.Body).Decode(v)
}

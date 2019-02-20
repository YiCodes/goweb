package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RouteMap map[string]string

func (m RouteMap) GetRoute(name string) string {
	route, ok := m[name]

	if ok {
		return route
	}

	return "/" + name
}

type HandlerSetupOptions struct {
	Route           RouteMap
	Codec           Codec
	OnRequestError  func(request *http.Request, response http.ResponseWriter, err error)
	OnResponseError func(request *http.Request, err error)
}

func defaultRequestErrorHandler(request *http.Request, response http.ResponseWriter, err error) {

}
func defaultResponseErrorHandler(request *http.Request, err error) {}

func NewHandleSetupOptions() *HandlerSetupOptions {
	opts := &HandlerSetupOptions{}
	opts.Codec = &JSONCodec{}
	opts.Route = make(RouteMap)
	opts.OnRequestError = defaultRequestErrorHandler
	opts.OnResponseError = defaultResponseErrorHandler

	return opts
}

func PostJson(url string, v interface{}) (*http.Response, error) {
	body, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewReader(body))
}

func ReadAsJson(reader io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(v)
}

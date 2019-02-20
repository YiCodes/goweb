package handlers

import (
	"github.com/YiCodes/goweb/web"
	"net/http"
	"github.com/YiCodes/goweb/example/hellowebapi/model"
)

func ConfigServeMuxHandler(mux *http.ServeMux, opts *web.HandlerSetupOptions) {
	routeMap := opts.Route
	msgCodec := opts.Codec
	onRequestError := opts.OnRequestError
	onResponseError := opts.OnResponseError
	mux.HandleFunc(routeMap.GetRoute("Hello"), func(w http.ResponseWriter, r *http.Request) {
		var a0 string
		var err error
		err = msgCodec.Decode(r, &a0)
		if err != nil {
			onRequestError(r, w, err)
			return
		}
		r0 := Hello(a0)
		err = msgCodec.Encode(w, &r0)
		if err != nil {
			onResponseError(r, err)
			return
		}
	})
	mux.HandleFunc(routeMap.GetRoute("AddUser"), func(w http.ResponseWriter, r *http.Request) {
		var a0 = new(model.User)
		var err error
		err = msgCodec.Decode(r, a0)
		if err != nil {
			onRequestError(r, w, err)
			return
		}
		AddUser(a0)
	})
	mux.HandleFunc(routeMap.GetRoute("GetUser"), func(w http.ResponseWriter, r *http.Request) {
		var a0 string
		var err error
		err = msgCodec.Decode(r, &a0)
		if err != nil {
			onRequestError(r, w, err)
			return
		}
		r0 := GetUser(a0)
		err = msgCodec.Encode(w, r0)
		if err != nil {
			onResponseError(r, err)
			return
		}
	})
	mux.HandleFunc(routeMap.GetRoute("AddOrder"), func(w http.ResponseWriter, r *http.Request) {
		var a0 = new(Order)
		var err error
		err = msgCodec.Decode(r, a0)
		if err != nil {
			onRequestError(r, w, err)
			return
		}
		AddOrder(a0)
	})
	mux.HandleFunc(routeMap.GetRoute("CustomHandle"), func(w http.ResponseWriter, r *http.Request) {
		CustomHandle(w)
	})
	mux.HandleFunc(routeMap.GetRoute("CustomHandle2"), func(w http.ResponseWriter, r *http.Request) {
		var a0 string
		var err error
		err = msgCodec.Decode(r, &a0)
		if err != nil {
			onRequestError(r, w, err)
			return
		}
		CustomHandle2(a0, r, w)
	})
}

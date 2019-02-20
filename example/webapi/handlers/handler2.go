package handlers

import "net/http"

type Order struct {
}

func CustomHandle(w http.ResponseWriter) {
	w.Write([]byte("hello world."))
}

func CustomHandle2(msg string, r *http.Request, w http.ResponseWriter) {
}

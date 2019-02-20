package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YiCodes/goweb/example/hellowebapi/model"

	"github.com/YiCodes/goweb/example/hellowebapi/handlers"
	"github.com/YiCodes/goweb/web"
)

func client() {
	time.Sleep(time.Second * 10)

	response, err := web.PostJson("http://127.0.0.1:8080/GetUser", "peter")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	user := model.User{}

	web.ReadAsJson(response.Body, &user)

	fmt.Println(user.Name)
}

func main() {
	mux := http.NewServeMux()

	handlerOpts := web.NewHandleSetupOptions()
	handlerOpts.Route["EmptyRequest"] = "/"

	handlers.ConfigServeMuxHandler(mux, handlerOpts)

	go client()

	log.Fatal(http.ListenAndServe(":8080", mux))
}

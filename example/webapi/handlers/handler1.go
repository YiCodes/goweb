package handlers

import (
	"fmt"

	"github.com/YiCodes/goweb/example/hellowebapi/model"
)

func Hello(name string) string {
	return fmt.Sprintf("hello %v", name)
}

func AddUser(user *model.User) {

}

func GetUser(name string) *model.User {
	return &model.User{Name: name}
}

func AddOrder(order *Order) {
}

package routes

import "net/http"

type Controller struct {
	Name string
}

type IController interface {
	GetPath() string
	Handler(res http.ResponseWriter, req *http.Request)
}

type Response struct {
}

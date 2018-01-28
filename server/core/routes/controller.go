package routes

import "net/http"

/**
* A Controller consists of a route and a payload.
* A Controller's handler determines what is returned to the user.
 */
type Controller struct {
	Name string
}

type IController interface {
	GetPath() string
	Handler(res http.ResponseWriter, req *http.Request)
}

type IControllerFactory interface {
}

type Response struct {
}

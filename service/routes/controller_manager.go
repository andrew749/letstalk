package routes

import (
	"net/http"
)

func RegisterControllers() {
	controllers := []IController{
		TestController{},
	}
	for _, controller := range controllers {
		http.HandleFunc(
			controller.GetPath(),
			controller.Handler,
		)
	}
}

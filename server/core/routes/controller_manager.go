package routes

import (
	"net/http"
)

func RegisterControllers() {
	controllers := []IController{
		TestController{},
		// ADD NEW CONTROLLERS HERE
	}
	for _, controller := range controllers {
		http.HandleFunc(
			controller.GetPath(),
			controller.Handler,
		)
	}
}

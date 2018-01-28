package routes

import (
	"log"
	"net/http"
)

func RegisterControllers() {

	controllers := []IController{
		TestController{},
		LoginController{},
		LoginCallbackController{},
		// ADD NEW CONTROLLERS HERE
	}
	log.Println("Registering controllers")
	for _, controller := range controllers {
		http.HandleFunc(
			controller.GetPath(),
			controller.Handler,
		)
	}
	log.Println("Registered controllers")
}

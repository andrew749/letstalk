package service

import (
	"net/http"

	"uwletstalk/service/routes"
)

func main() {
	routes.RegisterControllers()
	// start server
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"net/http"

	routes "uwletstalk/service/routes"
)

func main() {
	routes.RegisterControllers()
	// start server
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"net/http"

	"letstalk/server/core/routes"
)

func main() {
	routes.RegisterControllers()
	// start server
	http.ListenAndServe(":8080", nil)
}

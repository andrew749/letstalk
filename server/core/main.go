package main

import (
	"net/http"

	"letstalk/server/core/routes"
	"letstalk/server/core/secrets"
)

func main() {
	routes.RegisterControllers()
	secrets.GetSecrets()
	// start server
	http.ListenAndServe(":8080", nil)
}

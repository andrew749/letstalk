package main

import (
	"net/http"

	routes "./routes"
)

func main() {
	routes.RegisterControllers()
	// start server
	http.ListenAndServe(":8080", nil)
}

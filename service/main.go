package service

import (
	"net/http"

	routes "./routes"
)

func main() {
	_ := routes.RegisterControllers()
	// start server
	http.ListenAndServe(":8080", nil)
}

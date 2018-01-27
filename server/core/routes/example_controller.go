package routes

import (
	"net/http"
)

type TestController struct {
	Controller // quasi inheritance
}

func (c TestController) GetPath() string {
	return "/test"
}

func (c TestController) Handler(
	res http.ResponseWriter,
	req *http.Request,
) {
	res.Write([]byte("test controller"))
}

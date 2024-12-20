package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/topten1222/hello_sekai/pkg/response"
)

type healthCheck struct {
	App    string `json:"app"`
	Status string `json:"status"`
}

func (s *server) healthCheckService(c echo.Context) error {
	fmt.Println("343434")
	return response.SuccessResponse(c, http.StatusOK, &healthCheck{
		App:    s.cfg.App.Name,
		Status: "OK",
	})
}

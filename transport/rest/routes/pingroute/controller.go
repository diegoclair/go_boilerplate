package pingroute

import (
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
}

func NewController() *Controller {
	once.Do(func() {
		instance = &Controller{}
	})
	return instance
}

type pingResponse struct {
	Message string `json:"message"`
}

func (s *Controller) handlePing(c echo.Context) error {
	response := pingResponse{
		Message: "pong",
	}

	return c.JSON(200, response)
}

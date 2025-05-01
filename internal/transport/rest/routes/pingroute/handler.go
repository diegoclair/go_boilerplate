package pingroute

import (
	"net/http"
	"sync"

	echo "github.com/labstack/echo/v4"
)

var (
	instance *Handler
	once     sync.Once
)

type Handler struct {
}

func NewHandler() *Handler {
	once.Do(func() {
		instance = &Handler{}
	})
	return instance
}

type pingResponse struct {
	Message string `json:"message"`
}

func (s *Handler) handlePing(c echo.Context) error {
	response := pingResponse{
		Message: "pong",
	}

	return c.JSON(http.StatusOK, response)
}

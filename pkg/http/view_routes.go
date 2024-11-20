package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.Static("/", "templates/")
	e.File("/", "templates/index.html")
	e.File("/login", "templates/auth-login.html")
	e.File("/register", "templates/auth-register.html")
	e.GET("/sayhello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"hello": "world"})
	})
}

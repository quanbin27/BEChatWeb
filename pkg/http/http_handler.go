package http

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type HttpHandler struct {
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}
func (h *HttpHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/hello", h.SayHello)
}
func (h *HttpHandler) SayHello(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello world")
}
func NewGRPCClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

package http

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type HttpServer struct {
	addr string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{addr: addr}
}
func (s *HttpServer) Run() error {
	e := echo.New()
	conn := NewGRPCClient(":9000")
	defer conn.Close()
	subrouter := e.Group("/api/v1")
	httpHandler := NewHttpHandler(conn)
	httpHandler.RegisterRoutes(subrouter)
	log.Println("Listening on: ", s.addr)
	return http.ListenAndServe(s.addr, e)
}

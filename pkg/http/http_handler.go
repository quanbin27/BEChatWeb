package http

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"github.com/quanbin27/gRPC-Web-Chat/services/users/types"
	"github.com/quanbin27/gRPC-Web-Chat/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

type HttpHandler struct {
	grpcClient *grpc.ClientConn
}

func NewHttpHandler(grpcClient *grpc.ClientConn) *HttpHandler {
	return &HttpHandler{grpcClient: grpcClient}
}
func (h *HttpHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/hello", h.SayHello)
	e.POST("/register", h.RegisterHandler)
}
func (h *HttpHandler) SayHello(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello world")
}
func (h *HttpHandler) RegisterHandler(c echo.Context) error {
	var payload types.RegisterUserPayLoad
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors[0].Error()})
	}
	userClient := users.NewUserServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := userClient.Register(ctx, &users.RegisterRequest{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, res)
}
func NewGRPCClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

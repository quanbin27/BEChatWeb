package main

import (
	"github.com/quanbin27/gRPC-Web-Chat/pkg/http"
	"log"
)

func main() {
	httpServer := http.NewHttpServer(":1000")
	if err := httpServer.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"net/http"

	ServerConfig "github.com/guoxingx/fabtreehole/config/server"
	"github.com/guoxingx/fabtreehole/router"
)

func main() {
	r := router.InitRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  ServerConfig.ReadTimeout,
		WriteTimeout: ServerConfig.WriteTimeout,
	}

	server.ListenAndServe()
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/guoxingx/fabtreehole/config"
	"github.com/guoxingx/fabtreehole/pkg/fabconn"
	"github.com/guoxingx/fabtreehole/router"
)

func main() {
	r := router.InitRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  config.ServerConfig.ReadTimeout,
		WriteTimeout: config.ServerConfig.WriteTimeout,
	}

	err := fabconn.Setup()
	if err != nil {
		log.Panic(err)
	}

	server.ListenAndServe()
}

package server

import (
	"time"
)

const (
	RunMode      = "debug"
	Port         = 8000
	ReadTimeout  = time.Duration(10) * time.Second
	WriteTimeout = time.Duration(10) * time.Second
)

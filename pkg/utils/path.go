package utils

import (
	"os"
)

func GetProjectRoot() string {
	return os.Getenv("GOPATH") + "/src/github.com/guoxingx/fabtreehole"
}

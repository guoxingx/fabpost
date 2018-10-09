package config

import (
	"os"
	"time"

	"github.com/guoxingx/fabtreehole/pkg/utils"
)

const (
	RunMode      = "debug"
	Port         = 8000
	ReadTimeout  = time.Duration(10) * time.Second
	WriteTimeout = time.Duration(10) * time.Second
)

type Server struct {
	RunMode      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerConfig = Server{
	RunMode:      "debug",
	Port:         8000,
	ReadTimeout:  time.Duration(10) * time.Second,
	WriteTimeout: time.Duration(10) * time.Second,
}

type Fabric struct {
	ConfigFile      string
	OrdererID       string
	ChannelID       string
	ChaincodeID     string
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	UserName        string
	Members         []string
	Version         string
	Peers           []string
}

var FabricConfig = Fabric{
	ConfigFile:      utils.GetProjectRoot() + "/config/fabric-conn-default.yaml",
	OrdererID:       "orderer.example.com",
	ChannelID:       "mychannel",
	ChaincodeID:     "fabtreehole",
	ChannelConfig:   utils.GetProjectRoot() + "/network/default/config/channel.tx",
	ChaincodeGoPath: os.Getenv("GOPATH"),
	ChaincodePath:   "github.com/guoxingx/fabtreehole/chaincode",
	OrgAdmin:        "Admin",
	OrgName:         "org1",
	UserName:        "Admin",
	Members:         []string{"org1.example.com"},
	Version:         "0.1",
	Peers:           []string{"peer0.org1.example.com"},
}

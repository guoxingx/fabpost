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

type server struct {
	RunMode      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerConfig = server{
	RunMode:      "debug",
	Port:         8000,
	ReadTimeout:  time.Duration(10) * time.Second,
	WriteTimeout: time.Duration(10) * time.Second,
}

type fabric struct {
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

var FabricConfig = fabric{
	ConfigFile:      utils.GetProjectRoot() + "/config/fabric-conn-default.yaml",
	OrdererID:       "orderer.example.com",
	ChannelID:       "mychannel",
	ChaincodeID:     "fabtreehole",
	ChannelConfig:   utils.GetProjectRoot() + "/network/default/config/channel.tx",
	ChaincodeGoPath: os.Getenv("GOPATH"),
	ChaincodePath:   "github.com/guoxingx/fabtreehole/chaincode",
	OrgAdmin:        "Admin",
	OrgName:         "org1",
	UserName:        "User1",
	Members:         []string{"org1.dev.isu.com"},
	Version:         "0.1",
	Peers:           []string{"peer0.org1.dev.isu.com", "peer1.org1.dev.isu.com"},
}

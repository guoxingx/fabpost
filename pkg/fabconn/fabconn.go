package fabconn

import (
	"errors"
	"fmt"

	"github.com/guoxingx/fabtreehole/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	sdkconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type FabConn struct {
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
	sdk             *fabsdk.FabricSDK
	admin           *resmgmt.Client
	client          *channel.Client
	event           *event.Client
}

var conn *FabConn

func Setup() error {
	fmt.Println(config.FabricConfig.ConfigFile, event.Client{}, resmgmt.Client{})

	sdk, err := fabsdk.New(sdkconfig.FromFile(config.FabricConfig.ConfigFile))
	if err != nil {
		return errors.New(fmt.Sprintln("failed to create SDK: ", err.Error()))
	}
	fmt.Println("SDK created")
	fmt.Println(sdk)

	return nil
}

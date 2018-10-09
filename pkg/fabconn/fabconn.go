package fabconn

import (
	"errors"
	"fmt"
	"log"

	"github.com/guoxingx/fabtreehole/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	sdkconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

type FabConn struct {
	config.Fabric
	sdk    *fabsdk.FabricSDK
	admin  *resmgmt.Client
	client *channel.Client
	event  *event.Client
}

var conn *FabConn

type CCState int

const (
	CCNormal CCState = iota
	CCNotUpdated
	CCNotInstantiated
)

// Setup fabric go sdk
// Must be called firstly when using fabconn pkg
func Setup() error {
	conn := &FabConn{config.FabricConfig, nil, nil, nil, nil}

	sdk, err := fabsdk.New(sdkconfig.FromFile(conn.ConfigFile))
	if err != nil {
		return errors.New(fmt.Sprintln("failed to create SDK: ", err.Error()))
	}
	fmt.Println("SDK created")
	conn.sdk = sdk

	resourceManagerClientContext := sdk.Context(fabsdk.WithUser(conn.OrgAdmin), fabsdk.WithOrg(conn.OrgName))
	if err != nil {
		return errors.New(fmt.Sprintln("failed to load Admin identity: ", err.Error()))
	}

	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.New(fmt.Sprintln("failed to create channel management client from Admin identity: ", err.Error()))
	}
	conn.admin = resMgmtClient

	// create and join channel if not exit
	if !IsChannelExist() {
		err = CreateChannel()
		if err != nil {
			return errors.New(fmt.Sprintln("failed to create channel: ", err.Error()))
		}

		err = JoinChannel()
		if err != nil {
			return errors.New(fmt.Sprintln("failed to join channel: ", err.Error()))
		}
	}

	ccState := GetChaincodeState()
	// instantiate chaincode if not instantiated
	if ccState == CCNotInstantiated {
		// install chaincode if not installed
		if !IsChaincodeInstalled() {
			err = InstallChaincode()
			if err != nil {
				return errors.New(fmt.Sprintln("install chaincode failed: ", err.Error()))
			}
		}
		err = InstantiateChaincode()
		if err != nil {
			return errors.New(fmt.Sprintln("instantiate chaincode failed: ", err.Error()))
		}
	} else if ccState == CCNotUpdated {
		// update chaincode
		UpgradeChaincode()
	}

	clientContext := conn.sdk.ChannelContext(conn.ChannelID, fabsdk.WithUser(conn.UserName))
	client, err := channel.New(clientContext)
	if err != nil {
		return errors.New(fmt.Sprintln("channel client create failed: ", err.Error()))
	}
	fmt.Println("channel client created")
	conn.client = client

	event, err := event.New(clientContext)
	if err != nil {
		return errors.New(fmt.Sprintln("event client create failed: ", err.Error()))
	}
	conn.event = event

	return nil
}

func QueryChannels() (*pb.ChannelQueryResponse, error) {
	return conn.admin.QueryChannels(resmgmt.WithTargetEndpoints(conn.Peers[0]))
}

func IsChannelExist() bool {
	channelQueryResponse, err := QueryChannels()
	if err != nil {
		log.Panic("QueryChannels return error: %s", err)
	}

	// fmt.Println("channel query response: ", channelQueryResponse.Channels)
	for _, channel := range channelQueryResponse.Channels {
		if channel.ChannelId == conn.ChannelID {
			return true
		}
	}
	return false
}

func CreateChannel() error {
	mspClient, err := mspclient.New(conn.sdk.Context(), mspclient.WithOrg(conn.OrgName))
	if err != nil {
		return errors.New(fmt.Sprintln("failed to get admin signing identity: ", err.Error()))
	}

	adminIdentity, err := mspClient.GetSigningIdentity(conn.OrgAdmin)
	if err != nil {
		return errors.New(fmt.Sprintln("failed to get admin signing identity: ", err.Error()))
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         conn.ChannelID,
		ChannelConfigPath: conn.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}

	txID, err := conn.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(conn.OrdererID))
	if err != nil {
		return errors.New(fmt.Sprintln("failed to save channel: ", err.Error()))
	}
	if txID.TransactionID == "" {
		return errors.New(fmt.Sprintln("save channel response TransactionID empty: ", err.Error()))
	}
	fmt.Println("Channel created")
	return nil
}

func JoinChannel() error {
	err := conn.admin.JoinChannel(
		conn.ChannelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(conn.OrdererID))
	if err != nil {
		return errors.New(fmt.Sprintln("failed to make admin join channel: ", err.Error()))
	}
	fmt.Println("Channel joined: ", conn.ChannelID)
	return nil
}

func QueryInstalledChaincodes() (*pb.ChaincodeQueryResponse, error) {
	return conn.admin.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(conn.Peers[0]))
}

func QueryInstantiatedChaincodes() (*pb.ChaincodeQueryResponse, error) {
	return conn.admin.QueryInstantiatedChaincodes(conn.ChannelID, resmgmt.WithTargetEndpoints(conn.Peers[0]))
}

func IsChaincodeInstalled() bool {
	chaincodeQueryResponse, err := QueryInstalledChaincodes()
	if err != nil {
		log.Panic("QueryInstalledChaincodes return error: %s", err)
	}

	// fmt.Println("installed chaincode query response: ", chaincodeQueryResponse.Chaincodes)
	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		if chaincode.Name == conn.ChaincodeID {
			return true
		}
	}
	return false
}

func GetChaincodeState() CCState {
	chaincodeQueryResponse, err := QueryInstantiatedChaincodes()
	if err != nil {
		log.Panic("QueryInstantiatedChaincodes return error: %s", err)
	}

	// fmt.Println("instantiated chaincode query response: ", chaincodeQueryResponse.Chaincodes)
	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		if chaincode.Name == conn.ChaincodeID {
			if chaincode.Version == conn.Version {
				return CCNormal
			}
			return CCNotUpdated
		}
	}
	return CCNotInstantiated
}

func InstallChaincode() error {
	ccPkg, err := packager.NewCCPackage(conn.ChaincodePath, conn.ChaincodeGoPath)
	if err != nil {
		return errors.New(fmt.Sprintln("failed to create chaincode package: ", err.Error()))
	}
	fmt.Println("chaincode package created")

	installCCReq := resmgmt.InstallCCRequest{
		Name:    conn.ChaincodeID,
		Path:    conn.ChaincodePath,
		Version: conn.Version,
		Package: ccPkg}
	_, err = conn.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.New(fmt.Sprintln("chaincode install failed: ", err.Error()))
	}
	fmt.Println("chaincode installed")
	return nil
}

func InstantiateChaincode() error {
	ccPolicy := cauthdsl.SignedByAnyMember(conn.Members)

	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    conn.ChaincodeID,
		Path:    conn.ChaincodeGoPath,
		Version: conn.Version,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy}

	resp, err := conn.admin.InstantiateCC(conn.ChannelID, instantiateCCReq)
	if err != nil {
		return errors.New(fmt.Sprintln("chaincode instantiate failed: ", err.Error()))
	}
	if resp.TransactionID == "" {
		return errors.New(fmt.Sprintln("chaincode instantiate response TransactionID empty: ", err.Error()))
	}
	fmt.Println("chaincode instantiated")
	return nil
}

func UpgradeChaincode() error {
	ccPolicy := cauthdsl.SignedByAnyMember(conn.Members)

	upgradeCCReq := resmgmt.UpgradeCCRequest{
		Name:    conn.ChaincodeID,
		Version: conn.Version,
		Path:    conn.ChaincodeGoPath,
		Policy:  ccPolicy}

	resp, err := conn.admin.UpgradeCC(conn.ChannelID, upgradeCCReq)
	if err != nil {
		return errors.New(fmt.Sprintln("chaincode upgrade failed: ", err.Error()))
	}
	if resp.TransactionID == "" {
		return errors.New(fmt.Sprintln("chaincode upgrade response TransactionID empty: ", err.Error()))
	}
	fmt.Printf("chaincode upgrade to version %v success\n", conn.Version)
	return nil
}

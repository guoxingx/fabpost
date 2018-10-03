package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Rating struct{}

func (a *Rating) dispatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

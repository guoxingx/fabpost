package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Transaction struct{}

func (a *Transaction) dispatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func (a *Transaction) transact(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

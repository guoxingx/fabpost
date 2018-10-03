package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Account struct{}

// Generate message id like m-TxId-0
func getOwnerId(stub shim.ChaincodeStubInterface, index int) string {
	return fmt.Sprintf("o-%s-%d", stub.GetTxID(), index)
}

//
func (a *Account) dispatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	function, args, err := detachArgs(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	if function == "get" {
		return a.get(stub, args)
	} else if function == "add" {
		return a.add(stub, args)
	}

	return shim.Success(nil)
}

//
func (a *Account) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting account.get")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err := sanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get owner - " + args[0])
	}
	return shim.Success(valAsbytes)
}

// Create an Owner
//
// Inputs - []string
// 1. username
func (a *Account) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting account.add")

	if len(args) == 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err := sanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var owner Owner
	owner.ObjectType = OwnerObejctType
	owner.Id = getOwnerId(stub, 0)
	owner.Username = args[0]

	ownerAsByte, _ := json.Marshal(owner)
	err = stub.PutState(owner.Id, ownerAsByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

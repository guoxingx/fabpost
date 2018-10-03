package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Message struct{}

func (m *Message) dispatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	function, args, err := detachArgs(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	if function == "get" {
		return m.get(stub, args)
	} else if function == "add" {
		return m.add(stub, args)
	}

	return shim.Success(nil)
}

// Generate message id like m-TxId-0
func getMessageId(stub shim.ChaincodeStubInterface, index int) string {
	return fmt.Sprintf("m-%s-%d", stub.GetTxID(), index)
}

// get message by message id
func (m *Message) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var secretId, jsonResp string
	var err error
	fmt.Println("starting message.get")

	err = sanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	secretId = args[0]
	valAsbytes, err := stub.GetState(secretId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + secretId + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes)
}

// Create a secret
//
// Inputs - []string
// 1. content
// 2. owner id
func (m *Message) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting message.add")

	if len(args) == 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	err := sanitizeArguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	owner, err := getOwner(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	var secret Secret
	secret.ObjectType = SecretObjectType
	secret.Id = getMessageId(stub, 0)
	secret.Content = args[0]
	secret.Owner = OwnerRelation{owner.Id, owner.Username}

	secretAsByte, _ := json.Marshal(secret)
	err = stub.PutState(owner.Id, secretAsByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

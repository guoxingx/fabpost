package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type FabTreeHole struct{}

type Moduler interface {
	dispatch(stub shim.ChaincodeStubInterface, args []string) pb.Response
}

func main() {
	err := shim.Start(new(FabTreeHole))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

func (fm *FabTreeHole) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("FabTreeHole Is Starting Up")
	return shim.Success(nil)
}

func (fm *FabTreeHole) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	module, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + module)

	var m Moduler
	if module == "stuff" {
		m = &Message{}
	} else if module == "transaction" {
		m = &Transaction{}
	} else if module == "account" {
		m = &Account{}
	} else if module == "rating" {
		m = &Rating{}
	}
	return dispatchFunction(m, stub, args)
}

func dispatchFunction(m Moduler, stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return m.dispatch(stub, args)
}

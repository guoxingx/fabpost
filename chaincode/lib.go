package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// use first argument as function name
func detachArgs(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, errors.New("Incorrect number of arguments. Expecting 1 at least as function name.")
	}

	return args[0], args[1:], nil
}

// dumb input checking, look for empty strings
func sanitizeArguments(strs []string) error {
	for i, val := range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 32 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		}
	}
	return nil
}

func getOwner(stub shim.ChaincodeStubInterface, ownerId string) (Owner, error) {
	var owner Owner
	ownerAsByte, err := stub.GetState(ownerId)
	if err != nil {
		return owner, errors.New("Failed to get owner - " + ownerId)
	}
	json.Unmarshal(ownerAsByte, &owner)
	return owner, nil
}

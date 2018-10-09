package fabconn

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type RawRequest struct {
	Fcn          string
	Args         [][]byte
	TransientMap map[string][]byte
	EventID      string
}

// Query function
func Call(req RawRequest) ([]byte, error) {
	response, err := conn.client.Query(
		channel.Request{
			ChaincodeID:  conn.ChaincodeID,
			Fcn:          req.Fcn,
			Args:         req.Args,
			TransientMap: req.TransientMap,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to call: %v", err)
	}
	return response.Payload, nil
}

// Invoke function with write action
func Transact(req RawRequest) (string, error) {
	request := channel.Request{
		ChaincodeID:  conn.ChaincodeID,
		Fcn:          req.Fcn,
		Args:         req.Args,
		TransientMap: req.TransientMap,
	}
	if req.EventID != "" {
		return sendTransactWithEvent(request, req.EventID)
	}
	return sendTransact(request)
}

// send request to fabric
func sendTransact(request channel.Request) (string, error) {
	response, err := conn.client.Execute(request)
	if err != nil {
		return "", fmt.Errorf("transact failed: %v", err)
	}

	return string(response.TransactionID), nil
}

// send request to fabric
// and start listening transact event
func sendTransactWithEvent(request channel.Request, eventID string) (string, error) {
	reg, notifier, err := conn.event.RegisterChaincodeEvent(conn.ChaincodeID, eventID)
	if err != nil {
		return "", err
	}
	defer conn.event.Unregister(reg)

	response, err := conn.client.Execute(request)
	if err != nil {
		return "", fmt.Errorf("transact failed: %v", err)
	}

	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	return string(response.TransactionID), nil
}

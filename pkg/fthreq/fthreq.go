package fthreq

import (
	"github.com/guoxingx/fabtreehole/pkg/fabconn"
)

// Request struct for fabtreehole chaincode
type Request struct {
	Module       string
	Fuction      string
	Args         [][]byte
	TransientMap map[string][]byte
	EventID      string
}

// Convert type Request to fabconn.RawRequest
func (req *Request) ToRawRequest() *fabconn.RawRequest {
	args := [][]byte{[]byte(req.Fuction)}
	args = append(args, req.Args...)
	rreq := fabconn.RawRequest{
		Fcn:          req.Module,
		Args:         args,
		TransientMap: req.TransientMap,
		EventID:      req.EventID,
	}
	return rreq
}

func TestQuery() fabconn.RawRequest {
	return &fabconn.Request{
		Module:  "test",
		Fuction: "Query",
	}.ToRawRequest()
}

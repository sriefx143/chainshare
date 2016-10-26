// userinq
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Inqtoken struct {
	Token      string `json:"token"`
	Dated      string `json:"dated"`
	Inquiredby string `json:"inquiredby"`
}

type Inqtoken2 struct {
	Token      string `json:"token"`
	Dated      string `json:"dated"`
	Inquiredby string `json:"inquiredby"`
	InqData    string `json:"inqdata"`
}

type IDInquiry struct {
}

func main() {
	err := shim.Start(new(IDInquiry))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

func (i *IDInquiry) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (i *IDInquiry) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if function == "" || function == "post" {
		return i.postInq(stub, args)
	} else if function == "compare" {
		return i.compareInq(stub, args)
	}

	return nil, nil
}
func (i *IDInquiry) postInq(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	posterkey := args[0] + "-" + "inq"
	tm := time.Now().Format(time.RFC3339)
	tk := Inqtoken2{Token: args[1], Dated: tm, Inquiredby: args[0], InqData: args[2]}
	jStr, _ := json.Marshal(tk)
	stub.PutState(posterkey, jStr)

	return nil, nil
}

//chaincode to chaincode call test
//if same guy wants to audit or id request done for someone else --on behalf of
//eg. bank A or Id owner give id request for bank B for transaction and bank B checks it in flow
func (i *IDInquiry) compareInq(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	poster := args[0]
	inquiredby := args[1]
	tokenVal := args[2]
	var tk Inqtoken2

	posterkey := poster + "-" + "inq"
	stateData, _ := stub.GetState(posterkey)

	_ = json.Unmarshal(stateData, &tk)

	if tk.Inquiredby == inquiredby || tk.Inquiredby == poster {
		if tokenVal == tk.Token {
			return []byte("true"), nil
		}
	}
	b := []byte("")
	e := errors.New("INVALID_REQUEST_COMPARE_ERROR")
	return b, e
}

// Query is our entry point for queries
func (i *IDInquiry) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//poster := args[0]
	inquiredby := args[1]

	posterkey := inquiredby + "-" + "inq"
	stateData, _ := stub.GetState(posterkey)

	return stateData, nil
}

// userinq
package main

import (
	"encoding/json"

	"fmt"
	//	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//later put posted by (member)
type acctToken struct {
	Token string `json:"token"`
	Dated string `json:"dated"`
}

type Acct struct {
}

func main() {
	err := shim.Start(new(Acct))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

func (i *Acct) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (i *Acct) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	return i.postacct(stub, args)

}

//each acct post updates the state in a day
func (i *Acct) postacct(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	posterkey := args[0] + "-" + "accts"
	tm := time.Now().Format("2006-01-02")
	tk := acctToken{Token: args[1], Dated: tm}
	jStr, _ := json.Marshal(tk)
	stub.PutState(posterkey, jStr)

	return nil, nil
}

// Query is our entry point for queries
//dont worry about users yet think all have full access, just give out stored
func (i *Acct) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	posterkey := args[0] + "-" + "accts"
	stateData, _ := stub.GetState(posterkey)

	return stateData, nil
}

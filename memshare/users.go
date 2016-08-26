/*
test program
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	//	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// UserAccount example simple Chaincode implementation
type UserAccount struct {
}

//this is a data structure acting as pointer
type userprofile struct {
	Id      string            `json:"id"`
	Kvalues map[string]string `json:"kvalues"`
}

//this struct will be used by consumer  allowing inquiry
type user struct {
	Id       string    `json:"id"`
	Udata    []Idproof `json:"dlData"`
	OFACdata []OFAC    `json:"ofacdata"`
}

type Idproof struct {
	DLdata    []UserDLData `json:"dldata"`
	SsnData   []SsData     `json:"ssns"`
	TelcoDate []telco      `json:"telcodata"`
	Commdata  []comdata    `json:"commdata"`
}

type comdata struct {
	Companyname string `json:"companyname"`
	Address     string `json:"addr"`
	Founddate   string `json:"founddate"`
	Feduciary   bool   `json:"feduciary"`
	Active      bool   `json:"active"`
}

type OFAC struct {
	Flagdate string `json:"flagdate"`
	Active   bool   `json:"active"`
}

type telco struct {
	Phone          string `json:"phone"`
	memberprovider string `json:"memprovider"`
	startdate      string `json:"startdate"`
	enddate        string `json:"enddate"`
	reportdate     string `json:"reportdate"`
}

type SsData struct {
	Ssn    string `json:"SSN"`
	Active bool   `json:"active"`
}

//this struct will be used by consumer entity to read from for inquiry requests
type UserDLData struct {
	DLNum     string `json:"dlnum"`
	DLState   string `json:"dlstate"`
	IssueDate string `json:"issuedate"`
	ExpDate   string `json:"expdate"`
	FullNm    string `json:"fullname"`
	Dob       string `json:"dob"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(UserAccount))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// Init resets all the things
func (t *UserAccount) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var userid = args[0]

	var store = userid + "-" + "profile"

	t.inituser(stub, store, args)

	var shastore = userid + "-sha-" + "profile"

	t.inituser(stub, shastore, args)

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *UserAccount) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if function == "write" {
		t.writeUserHash(stub, args[0], args[1:])
	} else if function == "init" {

		var userid = args[0]
		var store = userid + "-" + "profile"

		t.inituser(stub, store, args)

		var shastore = userid + "-sha-" + "profile"
		t.inituser(stub, shastore, args)
	}

	return nil, nil
}

//invoke init user
func (t *UserAccount) inituser(stub *shim.ChaincodeStub, store string, args []string) ([]byte, error) {

	stub.PutState(store, []byte("INIT"))

	return nil, nil
}

func (t *UserAccount) writeUserHash(stub *shim.ChaincodeStub, userid string, args []string) ([]byte, error) {

	var shastore = userid + "-sha-" + "profile"

	var up userprofile

	br, er := stub.GetState(shastore)
	if er == nil {
		er2 := json.Unmarshal(br, &up)
		if er2 != nil {
			var mscurrent map[string]string = up.Kvalues
			_ = len(mscurrent)
			var msnew = make(map[string]string)
			for _, e := range args {
				r := strings.Split(e, ":")
				msnew[r[0]] = r[1]
			}

			up.Kvalues = msnew
			brnew, _ := json.Marshal(up)
			stub.PutState(shastore, brnew)
		}
	}

	return nil, nil

}

func (t *UserAccount) writeUserT(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil

}

func (t *UserAccount) writeUserCom(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil

}

func (t *UserAccount) writeUserOf(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil

}

func (t *UserAccount) writeUserDl(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil

}

func (t *UserAccount) writeUserSs(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil

}

//remove data in chaincode then post with new state record
func (t *UserAccount) remove(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil

}

// Query is our entry point for queries
func (t *UserAccount) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if function == "read" {
		retval, err := t.read(stub, args)
		return retval, err
	}
	return nil, errors.New("no function called")
}

func (t *UserAccount) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var aboutUser = args[1]

	var shastore string
	shastore = aboutUser + "-sha-" + "profile"

	barray, er := stub.GetState(shastore)

	return barray, er
}

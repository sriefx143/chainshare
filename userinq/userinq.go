// idinquire
package main

import (
	"crypto/sha512"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"

	"io/ioutil"
	//	"net"
	"net/http"
	"net/url"
	"os"

	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IdInquiry struct {
	Ssn      string `json:"ssn"`
	Inquirer string `json:"inquirer"`
	On       string `json:"on"`
}

type user struct {
	Id       string  `json:"id"`
	Udata    Idproof `json:"dlData"`
	OFACdata []OFAC  `json:"ofacdata"`
}

type Idproof struct {
	DLdata   []UserDLData `json:"dldata"`
	SsnData  []SsData     `json:"ssns"`
	Telco    []Telcodata  `json:"telcodata"`
	Commdata []comdata    `json:"commdata"`
}

type comdata struct {
	Companyname string `json:"companyname"`
	Address     string `json:"addr"`
	Founddate   string `json:"founddate"`
	Feduciary   bool   `json:"feduciary"`
	Active      bool   `json:"active"`
}

//this struct will be used in Id store and allowing inquiry
type OFAC struct {
	Flagdate string `json:"flagdate"`
	Active   bool   `json:"active"`
}

//this struct will be used in Id store and allowing inquiry
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

type Telcodata struct {
	Connexuskey string `json:"connexuskey"`
	Ssn         string `json:"ssn"`
	Fullname    string `json:fullname`
	Address     string `json:"address"`
	Startdate   string `json:"startdate"`
	Phone       string `json:"phone"`
}

func main() {
	var ssn string
	flag.StringVar(&ssn, "ssn", "ssn", "")
	var inquirer string
	flag.StringVar(&inquirer, "inq", "inq", "")

	var comp string
	flag.StringVar(&comp, "comp", "comp", "")

	flag.Parse()

	fmt.Println("inquiry details:")
	fmt.Println(ssn)
	fmt.Println(inquirer)
	if comp == "Y" || comp == "y" {
		compareInq(ssn, inquirer)

	} else {
		inquire(ssn, inquirer)
	}

}

func inquire(ssn string, inquiredby string) {

	var u user

	if ssn != "ssn" {
		session, err := mgo.Dial("localhost")
		if err != nil {
			fmt.Println("Error unable to reach database" + err.Error())

		}
		defer session.Close()
		db := session.DB("identitydb")

		col := db.C("idinquiries")

		idInblock := IdInquiry{Ssn: ssn, Inquirer: inquiredby}

		barrId, _ := json.Marshal(idInblock)

		sha_512 := sha512.New()
		shaoutput := sha_512.Sum(barrId)
		hxString := hex.EncodeToString(shaoutput)

		fmt.Println(hxString)

		inquiredOn := time.Now().UTC().Format(time.RFC3339)

		idIn := IdInquiry{Ssn: ssn, Inquirer: inquiredby, On: inquiredOn}

		db.C("users").Find(bson.M{"udata.ssndata": bson.M{"$elemMatch": bson.M{"ssn": ssn}}}).One(&u)

		if u.Id != "" {
			fmt.Println("found user:")
			fmt.Println(u)
		} else {
			fmt.Println("user not found:")
		}

		ts := string(time.Now().UnixNano())

		col.Insert(&idIn)

		post2userinqbc(idIn, ts)

		fmt.Println("inquiry done on:")
		fmt.Println(idIn)
		fmt.Println("token:", ts)

	}
}

func compareInq(ssn string, inquirer string) bool {
	if ssn != "ssn" {
		var i IdInquiry
		session, err := mgo.Dial("localhost")
		if err != nil {
			fmt.Println("Error unable to reach database" + err.Error())

		}
		defer session.Close()
		db := session.DB("identitydb")
		db.C("idinquiries").Find(bson.M{"ssn": ssn}).One(&i)

		if i.Ssn == ssn && inquirer == i.Inquirer {
			fmt.Println("inquiry compare successful")
			return true
		}
	}

	fmt.Println("inquiry compare unsuccessful")
	return false
}

func post2userinqbc(inqId IdInquiry, ts string) {

	certfile, _ := ioutil.ReadFile("C:\\Sri\\projectsdocument\\BlockChain\\sep-cert.cer")
	cert := tls.Certificate{Certificate: [][]byte{certfile}}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	//t := &http.Transport{TLSClientConfig: &config}
	//client := &http.Client{Transport: t}

	var body string = "{" +
		"\"jsonrpc\": \"2.0\"," +
		"\"method\": \"invoke\"," +
		"\"params\": {" +
		"\"type\": 1," +
		"\"chaincodeID\": {" +
		"\"name\": %s" +
		"}," +
		"\"ctorMsg\": {" +
		"\"function\": %s," +
		"\"args\": [%s]}," +
		"\"secureContext\": \"%s\"" +
		"}," +
		"\"id\": 2 " +
		"}"

	sha := sha512.New()
	sha.Write([]byte(inqId.Ssn + inqId.Inquirer))
	op := sha.Sum(nil)

	hxstring := hex.EncodeToString(op)

	var argsstr string = "\"" + inqId.Inquirer + "\",\"" + hxstring + "\",\"" + ts + "\""

	//chainid := "\"bc35f431991c39a1e255c663fcaa80a607c438f497d967ec5cefc2d2415442db89713a085d69e7df0663de585a318c1c1c00c4b1b89af82401f1abaf709b3a66\""
	chainid := "\"c13948d95150a7d51040f7cec06e5f520743361d9c754f5f1b147ab9caaf07ad192e979a66a3983e5d87b838baeb3ee9cef5a3fb0f14dcc31ced3b2af24abe2a\""
	body = fmt.Sprintf(body, chainid, "\"post\"", argsstr, "WebAppAdmin")
	fmt.Println(body)
	var netTransport = &http.Transport{TLSClientConfig: &config}
	//,

	//Dial: (&net.Dialer{
	//Timeout: 5 * time.Second,
	//}).Dial,
	//TLSHandshakeTimeout: 5 * time.Second,

	//var netClient = &http.Client{
	//Timeout:   time.Second * 100,
	//Transport: netTransport,
	//}

	//resp, callerr := netClient.Post("", "application/json", strings.NewReader(body))
	os.Setenv("HTTP_PROXY", "http://172.18.100.15:18717")
	client := &http.Client{Transport: netTransport}
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse("http://sxk143:AtmavaIdham@1@172.18.100.15:18717")
	netTransport.Proxy = http.ProxyURL(url_proxy)
	resp, callerr := client.Post("https://aff11732-f1ad-45f4-ae70-2daffd4492d6_vp1.us.blockchain.ibm.com:443/chaincode", "application/json", strings.NewReader(body))
	fmt.Println(resp)
	fmt.Println("Error---------------------------")
	fmt.Println(callerr)
	fmt.Println("End Error---------------------------")

	resbody, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resbody))
	fmt.Println(err)
}

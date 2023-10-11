package kademlia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTPInterface(t *testing.T) {
	//TEsting init
	kademlia, _ := InitJoin("172.30.80.1", 2001)
	go InitHTTPInterface(kademlia, kademlia.RoutingTable)
	testInput := "test123"
	data := []byte(testInput)
	//Seting up values
	kademlia.Store(data)

	//Testing Get
	res, err := http.Get("http://172.30.80.1:5000/objects/74657374313233da39a3ee5e6b4b0d3255bfef95")
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	fmt.Println("client: got response!, code: ", res)
	respString := ""
	json.Unmarshal(resBody, &respString)

	if respString != testInput {
		t.Errorf("Expected output: %s, got: %s", testInput, respString)
	}

	//Testing Post
	body, _ := json.Marshal("test456")

	res, err = http.Post("http://172.30.80.1:5000/objects/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	resBody, err = ioutil.ReadAll(res.Body)
	respString = ""
	json.Unmarshal(resBody, &respString)

	if respString != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Errorf("Expected output: %s, got: %s", "da39a3ee5e6b4b0d3255bfef95601890afd80709", respString)
	}

}

package kademlia

import (
	"testing"
)

func TestHTTPClientInputHandler(t *testing.T) {

	inputString := []string{"POST", "8080", "test1111"}
	//inputString := []string{"GET", "8084", "hash"}
	httpClientInputHandler(inputString)

	//httpClient()
}

package kademlia

import (
	"testing"
)

func TestHTTPClientInputHandler(t *testing.T) {

	//inputString := []string{"GET", "192.168.0.188", "test12345"}
	inputString := []string{"POST", "localhost:8084", "test12345"}
	httpClientInputHandler(inputString)
}

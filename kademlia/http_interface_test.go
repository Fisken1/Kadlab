package kademlia

import (
	"testing"
)

func TestHTTPInterface(t *testing.T) {
	kademlia := &Kademlia{
		RoutingTable: NewRoutingTable(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost", 8080)),
	}

	InitHTTPInterface(kademlia, kademlia.RoutingTable)

	select {}

}

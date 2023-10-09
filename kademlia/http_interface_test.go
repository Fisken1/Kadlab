package kademlia

import (
	"testing"
)

func TestHTTPInterface(t *testing.T) {
	kademlia := &Kademlia{
		RoutingTable: NewRoutingTable(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost", 8080)),
	}
	println("before start")
	InitHTTPInterface(kademlia, kademlia.RoutingTable)
	println("started")
	select {}
	println("Escaped")
}

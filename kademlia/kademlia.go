package kademlia

import (
	_ "fmt"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Hashmap      map[string][]byte
}

func InitKademila(me Contact) *Kademlia {
	node := &Kademlia{
		RoutingTable: NewRoutingTable(me),
		Hashmap:      make(map[string][]byte),
	}
	return node
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

package kademlia

import (
	"fmt"
	"testing"
)

func TestInitNode(t *testing.T) {
	id1 := NewRandomKademliaID()
	id2 := NewRandomKademliaID()
	if id1 == id2 {
		t.Error("THEY ARE THE SAME")
	}
	fmt.Println(id1.String(), id2.String())
}

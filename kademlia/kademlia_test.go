package kademlia

import (
	"fmt"
	"strings"
	"testing"
)

var bootstrap *Kademlia = nil

func TestInitNode(t *testing.T) {
	// Initialize a Kademlia node and verify its fields are set correctly.
	me := NewContact(NewRandomKademliaID(), "127.0.0.1", 12345)
	node := InitNode(me)

	// Check if the RoutingTable, Hashmap, and other fields are correctly initialized.
	if node.RoutingTable == nil {
		t.Error("RoutingTable is not initialized")
	}

	if node.storagehandler == nil {
		t.Error("Hashmap is not initialized")
	}

	if node.alpha != 3 {
		t.Errorf("Expected alpha to be 3, got %d", node.alpha)
	}

	if node.k != 10 {
		t.Errorf("Expected k to be 4, got %d", node.k)
	}
	// Add more assertions as needed.
}

func TestInitJoinBootstrap(t *testing.T) {
	// Initialize a Kademlia bootstrap node and verify its fields.
	bootstrap, _ = InitJoin("130.240.65.81", 5000)

	// Check if the node is correctly initialized as a bootstrap node.
	if !bootstrap.bootstrap {
		t.Error("Expected bootstrap node, got non-bootstrap node")
	}

	// Add more assertions as needed.
}

func TestInitJoinNonBootstrap(t *testing.T) {
	// Initialize a Kademlia non-bootstrap node and verify its fields.
	//bootstrap, _ = InitJoin("192.168.1.26", 5000)
	nonBootstrapNode, _ := InitJoin("127.0.0.1", 2000)

	nonBootstrapNodeContacts := nonBootstrapNode.RoutingTable.FindClosestContacts(bootstrap.RoutingTable.me.ID, 1)
	bootstrapContacts := bootstrap.RoutingTable.FindClosestContacts(nonBootstrapNode.RoutingTable.me.ID, 1)

	// Check if the node is correctly initialized as a non-bootstrap node.
	if nonBootstrapNode.bootstrap {
		t.Error("Expected non-bootstrap node, got bootstrap node")
	}

	if nonBootstrapNodeContacts[0].ID.String() != bootstrap.RoutingTable.me.ID.String() {
		t.Error("Expected", bootstrap.RoutingTable.me.ID)
	}

	if bootstrapContacts[0].ID.String() != nonBootstrapNode.RoutingTable.me.ID.String() {
		t.Error("Expected", nonBootstrapNode.RoutingTable.me.ID.String())
	}
}

func TestGetBootstrapIP(t *testing.T) {
	// Test the GetBootstrapIP function with different IP inputs.
	bootstrapIP := GetBootstrapIP("130.240.65.81")
	if bootstrapIP != "130.240.65.81" {
		t.Errorf("Expected bootstrap IP '130.240.65.81', got %s", bootstrapIP)
	}

}

func TestStore(t *testing.T) {
	nonBootstrapNode1, _ := InitJoin("127.0.0.1", 2001)

	// Test the Store function by storing data and verifying its presence.
	testData := []byte("test data")

	// Run the Store method
	result, _ := nonBootstrapNode1.Store(testData)
	fmt.Println(result, " nonBootstrapNode1", nonBootstrapNode1.RoutingTable.me.String())
	if result == "" {
		t.Error("Store did not return a valid hash.")
	}
}

func TestLookupData(t *testing.T) {
	nonBootstrapNode1, _ := InitJoin("127.0.0.1", 2004)

	// Run the Store method
	result, con, _, err := nonBootstrapNode1.LookupData("746573742064617461da39a3ee5e6b4b0d3255bfef95601890afd80709")
	fmt.Println(result, " nonBootstrapNode1", nonBootstrapNode1.RoutingTable.me.String())
	if result[0].String() == "" {
		t.Error("Store did not return a valid hash.")
	}
	if con.String() == "" {
		t.Error("Store did not return a valid hash.")
	}

	if err != nil {
		t.Error("Store did not return a valid hash.")
	}

	input := []string{"get", "746573742064617461da39a3ee5e6b4b0d3255bfef95601890afd80709"}
	outputGet := CliHandler(input, nonBootstrapNode1)

	if !strings.Contains(outputGet, "Found data") {
		t.Errorf("Expected output: %s, got: %s", "Found data: test data from contact: e36726ff43292663c457f3c5692130835537a98a. At adress: 127.0.0.1", outputGet)
	}
}

func TestBucketLen(t *testing.T) {
	node := &Kademlia{
		RoutingTable:   NewRoutingTable(NewContact(NewKademliaID(NewRandomKademliaID().String()), "127.0.0.1", 2005)),
		storagehandler: &StorageHandler{},
		alpha:          3,
		k:              10,
	}
	node.RoutingTable.buckets[10].AddContact(bootstrap.RoutingTable.me)
	result := node.RoutingTable.buckets[10].Len()
	fmt.Println("Len of bucket 1 ", result)
	if result != 1 {
		t.Error("Expected 2 nodes in bucket 1 due to prev tests")
	}
}

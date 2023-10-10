package kademlia

import (
	"fmt"
	"strings"
	"testing"
)

var bootstrap *Kademlia = nil
var storagenode *Kademlia = nil

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
	bootstrap, _ = InitJoin("130.240.109.105", 5000)

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
	bootstrapIP := GetBootstrapIP("130.240.109.105")
	if bootstrapIP != "130.240.109.105" {
		t.Errorf("Expected bootstrap IP '130.240.65.81', got %s", bootstrapIP)
	}

}

func TestStore(t *testing.T) {
	storagenode, _ = InitJoin("127.0.0.1", 2001)

	// Test the Store function by storing data and verifying its presence.
	testData := []byte("test data")

	// Run the Store method
	result, _ := storagenode.Store(testData)
	fmt.Println(result, " nonBootstrapNode1", storagenode.RoutingTable.me.String())
	if result == "" {
		t.Error("Store did not return a valid hash.")
	}

}

func TestLookupData(t *testing.T) {
	nonBootstrapNode2, _ := InitJoin("127.0.0.1", 2004)

	// Run the Store method
	result, con, _, err := nonBootstrapNode2.LookupData("746573742064617461da39a3ee5e6b4b0d3255bfef95601890afd80709")
	fmt.Println(result, " nonBootstrapNode1", nonBootstrapNode2.RoutingTable.me.String())
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
	outputGet := CliHandler(input, nonBootstrapNode2)

	if !strings.Contains(outputGet, "Found data") {
		t.Errorf("Expected output: %s, got: %s", "Found data: test data from contact: e36726ff43292663c457f3c5692130835537a98a. At adress: 127.0.0.1", outputGet)
	}

	outputStored := storagenode.storagehandler.getStoredData()
	if !strings.Contains(string(outputStored[0].Value), "test data") {
		t.Errorf("Expected output: %s, got: %s", "test data", string(outputStored[0].Value))
	}
	outputStored2 := storagenode.storagehandler.getUploadedData()
	if !strings.Contains(outputStored2[0], "746573742064617461da39a3ee5e6b4b0d3255bf") {
		t.Errorf("Expected output: %s, got: %s", "746573742064617461da39a3ee5e6b4b0d3255bf", outputStored2[0])
	}

	outputBool := storagenode.ForgetData("746573742064617461da39a3ee5e6b4b0d3255bfef95601890afd80709")
	if outputBool {
		t.Error("BIG ERROR IN FORGET")
	}

	outputBool2 := storagenode.RoutingTable.removeContactAtFront(bootstrap.RoutingTable.me)
	if !outputBool2 {
		t.Error("BIG ERROR IN removeContactAtFront")

	}

	outputStored3 := storagenode.storagehandler.clearExpiredData()
	if outputStored3 == nil {
		t.Errorf("BIG ERROR IN CLEAR")
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

func TestGetAlphaContacts(t *testing.T) {
	// Create a Kademlia instance.

	nonBootstrapNode1, _ := InitJoin("127.0.0.1", 2020)

	// Create a map to track contacted nodes.
	contactedMap := make(map[string]bool)

	// Define the number of alpha contacts you want to retrieve.
	alpha := 3

	// Call the getAlphaContacts function.
	alphaContacts := nonBootstrapNode1.getAlphaContacts(&bootstrap.RoutingTable.me, alpha, contactedMap)

	// Check if the length of alphaContacts is equal to alpha.
	if len(alphaContacts) != alpha {
		t.Errorf("Expected %d alpha contacts, but got %d", alpha, len(alphaContacts))
	}

}

func TestGetClosestNode(t *testing.T) {
	// Create a Kademlia instance.

	targetID := NewRandomKademliaID()

	// Create a slice of contacts for testing.
	contacts := []Contact{
		NewContact(NewRandomKademliaID(), "127.0.0.1", 12345),
		NewContact(NewRandomKademliaID(), "127.0.0.1", 12346),
		NewContact(NewRandomKademliaID(), "127.0.0.1", 12347),
		NewContact(NewRandomKademliaID(), "127.0.0.1", 12348),
	}

	// Call the getClosestNode function.
	closest := storagenode.getClosestNode(*targetID, contacts)

	// Check if the closest contact is not nil.
	if closest == nil {
		t.Error("Expected a closest contact, but got nil")
	}

	// Add more specific assertions as needed based on your implementation and test requirements.
}

package kademlia

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestInitNode(t *testing.T) {
	mockContact := Contact{
		ID:      NewRandomKademliaID(),
		Address: "127.0.0.1:12345",
	}

	node := InitNode(mockContact)

	if node == nil {
		t.Error("InitNode should return a non-nil node")
	}

	if node.RoutingTable == nil {
		t.Error("RoutingTable should be initialized")
	}

	if node.Hashmap == nil {
		t.Error("Hashmap should be initialized")
	}

	if node.alpha != 3 {
		t.Errorf("Expected alpha to be 3, but got %d", node.alpha)
	}

	if node.net == nil {
		t.Error("net field should be set")
	}
}

func TestGetAlphaContacts(t *testing.T) {
	// Create a Kademlia instance for testing.
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:12345", 12345)
	kademlia := InitNode(me)

	// Create a list of contacts for testing.
	contacts := []Contact{
		NewContact(NewRandomKademliaID(), "127.0.0.1:1111", 1111),
		NewContact(NewRandomKademliaID(), "127.0.0.1:2222", 2222),
		NewContact(NewRandomKademliaID(), "127.0.0.1:3333", 3333),
		NewContact(NewRandomKademliaID(), "127.0.0.1:4444", 4444),
	}

	// Add these contacts to the routing table.
	for _, contact := range contacts {
		kademlia.RoutingTable.AddContact(contact)
	}

	// Call the function to get alpha contacts.
	alpha := 2
	alphaContacts := kademlia.getAlphaContacts(&contacts[0], []Contact{}, alpha, make(map[string]bool))

	// Check if the length of alphaContacts is equal to alpha.
	if len(alphaContacts) != alpha {
		t.Errorf("Expected %d alpha contacts, but got %d", alpha, len(alphaContacts))
	}

	// Check if the selected alpha contacts are in the list of available contacts.
	for _, contact := range alphaContacts {
		found := false
		for _, c := range contacts {
			if contact.ID.Equals(c.ID) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Selected contact is not in the list of available contacts")
		}
	}
}

func TestInitJoinAsBootstrapNode(t *testing.T) {
	ip := "130.240.64.25"
	ipBootstrap := "130.240.64.25:8080" // Same IP as the current node for bootstrap
	port := 8080

	node, err := InitJoin(ip, ipBootstrap, port)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Use a channel for synchronization.
	done := make(chan bool)

	// Start a goroutine to close the network and signal when done.
	go func() {
		// Sleep for a while to allow the network to start.
		time.Sleep(1 * time.Second)
		errorro := node.net.server.Close()
		if errorro != nil {
			t.Errorf("Error closing the server: %v", errorro)
		}
		done <- true
	}()

	// Wait for the goroutine to finish before proceeding with other checks.
	<-done

	if node == nil {
		t.Error("InitJoin should return a non-nil node")
	}

	if node.RoutingTable == nil {
		t.Error("RoutingTable should be initialized")
	}

	if node.Hashmap == nil {
		t.Error("Hashmap should be initialized")
	}

	if node.alpha != 3 {
		t.Errorf("Expected alpha to be 3, but got %d", node.alpha)
	}

	if node.net == nil {
		t.Error("net field should be set")
	}

	// Check that the node's contact information is set correctly.
	if node.RoutingTable.me.Address != ip+":"+strconv.Itoa(port) {
		t.Errorf("Expected node address to be %s, but got %s", ip, node.RoutingTable.me.Address)
	}

	if node.RoutingTable.me.Port != port {
		t.Errorf("Expected node port to be %d, but got %d", port, node.RoutingTable.me.Port)
	}

}

func TestInitJoinAsRegularNode(t *testing.T) {
	ip := "130.240.64.25"
	ipBootstrap := "130.240.64.25:8080" // Same IP as the current node for bootstrap
	port := 8080

	node, err := InitJoin(ip, ipBootstrap, port)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Use a channel for synchronization.
	done := make(chan bool)

	// Start a goroutine to close the network and signal when done.
	go func() {
		// Sleep for a while to allow the network to start.
		time.Sleep(1 * time.Second)
		errorro := node.net.server.Close()
		if errorro != nil {
			t.Errorf("Error closing the server: %v", errorro)
		}
		done <- true
	}()

	// Wait for the goroutine to finish before proceeding with other checks.
	<-done

	if node == nil {
		t.Error("InitJoin should return a non-nil node")
	}

	if node.RoutingTable == nil {
		t.Error("RoutingTable should be initialized")
	}

	if node.Hashmap == nil {
		t.Error("Hashmap should be initialized")
	}

	if node.alpha != 3 {
		t.Errorf("Expected alpha to be 3, but got %d", node.alpha)
	}

	if node.net == nil {
		t.Error("net field should be set")
	}

	// Check that the node's contact information is set correctly.
	if node.RoutingTable.me.Address != ip {
		t.Errorf("Expected node address to be %s, but got %s", ip, node.RoutingTable.me.Address)
	}

	if node.RoutingTable.me.Port != port {
		t.Errorf("Expected node port to be %d, but got %d", port, node.RoutingTable.me.Port)
	}
}

func TestGetClosestNode(t *testing.T) {
	// Create a Kademlia instance for testing.
	kademlia := &Kademlia{
		RoutingTable: NewRoutingTable(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost", 8080)),
		Hashmap:      make(map[string][]byte),
	}

	// Create a specific target ID.
	targetID := kademlia.RoutingTable.me.ID
	q := targetID.String()
	fmt.Print(q)
	// Create some queried contacts with known distances.
	queriedContacts := []Contact{
		NewContact(NewKademliaID("1212121200000000000000000000000000000000"), "130.240.64.89", 8080),
		NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "130.240.64.25", 8080),
		NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "130.240.64.26", 8080),
		NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "130.240.64.27", 8080),
		NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "130.240.64.28", 8080),
		NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "130.240.64.29", 8080),
	}

	// Calculate the XOR distances for the queried contacts.
	for i := range queriedContacts {
		queriedContacts[i].CalcDistance(targetID)
	}

	// Get the closest node.
	closestNode := kademlia.getClosestNode(*targetID, queriedContacts)

	// Add your assertions based on the known distances and expected results.
	// ...

	// Example assertions:
	if closestNode == nil {
		t.Error("Expected closestNode to be non-nil, but got nil")
	}
	if closestNode.ID != queriedContacts[0].ID {
		t.Errorf("Expected closestNode ID to be %s, but got %s", queriedContacts[0].ID.String(), closestNode.ID.String())
	}
	if closestNode.Address != queriedContacts[0].Address {
		t.Errorf("Expected closestNode Address to be %s, but got %s", queriedContacts[0].Address, closestNode.Address)
	}
	if closestNode.Port != queriedContacts[0].Port {
		t.Errorf("Expected closestNode Port to be %d, but got %d", queriedContacts[0].Port, closestNode.Port)
	}
}

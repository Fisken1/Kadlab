package kademlia

import (
	"fmt"
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
	ip := "192.168.1.25"
	bootstrap, err := InitJoin(ip, 57707)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Use a channel for synchronization.
	done := make(chan bool)

	// Start a goroutine to close the network and signal when done.
	go func() {
		// Sleep for a while to allow the network to start.
		time.Sleep(1 * time.Second)
		errorro := bootstrap.net.server.Close()
		if errorro != nil {
			t.Errorf("Error closing the server: %v", errorro)
		}
		done <- true
	}()

	// Wait for the goroutine to finish before proceeding with other checks.
	<-done

	if bootstrap == nil {
		t.Error("InitJoin should return a non-nil node")
	}

	if bootstrap.RoutingTable == nil {
		t.Error("RoutingTable should be initialized")
	}

	if bootstrap.Hashmap == nil {
		t.Error("Hashmap should be initialized")
	}

	if bootstrap.alpha != 3 {
		t.Errorf("Expected alpha to be 3, but got %d", bootstrap.alpha)
	}

	if bootstrap.net == nil {
		t.Error("net field should be set")
	}

}

// TODO FIX THIS TEST ITS BROKEN!!!!
func TestInitJoinAsRegularNode(t *testing.T) {
	ip := "192.168.1.26"
	port := 57707
	node, err := InitJoin(ip, port)
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

	correctnode := queriedContacts[2]

	// Get the closest node.
	closestNode := kademlia.getClosestNode(*targetID, queriedContacts)

	if closestNode == nil {
		t.Error("Expected closestNode to be non-nil, but got nil")
	}
	if closestNode.ID != correctnode.ID {
		t.Errorf("Expected closestNode ID to be %s, but got %s", queriedContacts[0].ID.String(), closestNode.ID.String())
	}
	if closestNode.Address != correctnode.Address {
		t.Errorf("Expected closestNode Address to be %s, but got %s", queriedContacts[0].Address, closestNode.Address)
	}
	if closestNode.Port != correctnode.Port {
		t.Errorf("Expected closestNode Port to be %d, but got %d", queriedContacts[0].Port, closestNode.Port)
	}
}
func TestContactLessFalse(t *testing.T) {
	// Create two Contact instances with the same distance
	contact1 := &Contact{
		distance: NewKademliaID("1111111200000000000000000000000000000000"),
	}
	contact2 := &Contact{
		distance: NewKademliaID("1111111200000000000000000000000000000000"),
	}

	// Compare 'distance' attributes using the Less method
	result := contact1.Less(contact2)

	// Expect the result to be false
	if result != false {
		t.Errorf("Expected contact1.distance to be less than contact2.distance, but got true")
	}
}

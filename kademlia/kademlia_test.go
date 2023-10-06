package kademlia

import (
	"testing"
)

func TestInitNode(t *testing.T) {
	// Initialize a Kademlia node and verify its fields are set correctly.
	me := NewContact(NewRandomKademliaID(), "127.0.0.1", 12345)
	node := InitNode(me)

	// Check if the RoutingTable, Hashmap, and other fields are correctly initialized.
	if node.RoutingTable == nil {
		t.Error("RoutingTable is not initialized")
	}

	if node.Hashmap == nil {
		t.Error("Hashmap is not initialized")
	}

	if node.alpha != 3 {
		t.Errorf("Expected alpha to be 3, got %d", node.alpha)
	}

	if node.k != 4 {
		t.Errorf("Expected k to be 4, got %d", node.k)
	}
	// Add more assertions as needed.
}

/*
	func TestInitJoinBootstrap(t *testing.T) {
		// Initialize a Kademlia bootstrap node and verify its fields.
		bootstrapNode, _ := InitJoin("172.168.0.2", 5000)

		// Check if the node is correctly initialized as a bootstrap node.
		if !bootstrapNode.bootstrap {
			t.Error("Expected bootstrap node, got non-bootstrap node")
		}

		// Verify that the bootstrapContact and network are set up as expected.
		if bootstrapNode.bootstrapContact == nil {
			t.Error("Bootstrap contact is not initialized")
		}

		// Add more assertions as needed.
	}

	func TestInitJoinNonBootstrap(t *testing.T) {
		// Initialize a Kademlia non-bootstrap node and verify its fields.
		_, _ = InitJoin("172.168.0.2", 5000)
		nonBootstrapNode, _ := InitJoin("172.168.0.3", 5000)

		// Check if the node is correctly initialized as a non-bootstrap node.
		if nonBootstrapNode.bootstrap {
			t.Error("Expected non-bootstrap node, got bootstrap node")
		}

		// Verify that the bootstrapContact and network are set up as expected.
		if nonBootstrapNode.bootstrapContact != nil {
			t.Error("Expected nil bootstrap contact for non-bootstrap node")
		}

		// Add more assertions as needed.
	}
*/
func TestGetBootstrapIP(t *testing.T) {
	// Test the GetBootstrapIP function with different IP inputs.
	bootstrapIP := GetBootstrapIP("192.168.1.1")
	if bootstrapIP != "172.168.0.2" {
		t.Errorf("Expected bootstrap IP '172.168.0.2', got %s", bootstrapIP)
	}

	// Add more test cases for different IP inputs.
}

//func TestStore(t *testing.T) {
// Test the Store function by storing data and verifying its presence.

//}

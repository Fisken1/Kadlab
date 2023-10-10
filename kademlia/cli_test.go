package kademlia

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestCliHandler(t *testing.T) {
	// Create a mock Kademlia node for testing

	me := NewContact(NewKademliaID(BootstrapKademliaID), "127.0.0.1", 9000)
	node := InitNode(me)

	input := []string{"printID"}
	outputID := CliHandler(input, node)

	me2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "127.0.0.1", 9001)
	node2 := InitNode(me2)
	node.RoutingTable.buckets[10].AddContact(me2)
	input = []string{"getContact"}
	outputContacts := CliHandler(input, node)

	input = []string{"put", "hej"}
	outputPut := CliHandler(input, node)

	input = []string{"NOTACOMMAND"}
	outputDef := CliHandler(input, node)

	if outputID != node.RoutingTable.me.ID.String() {
		t.Errorf("Expected output: %s, got: %s", node.RoutingTable.me.ID.String(), outputID)
	}
	if outputContacts != node2.RoutingTable.me.ID.String() {
		t.Errorf("Expected output: %s, got: %s", node2.RoutingTable.me.ID.String(), outputContacts)
	}
	if outputPut != "68656ada39a3ee5e6b4b0d3255bfef95601890af\n" {
		t.Errorf("Expected output: %s, got: %s", "68656ada39a3ee5e6b4b0d3255bfef95601890af", outputPut)
	}

	expectedUsage := "Usage: \n\tput [contents] \n\t\tTakes a single argument, the contents of the file you are uploading, and outputs the\n\t\thash of the object, if it could be uploaded successfully.\n\tget [hash] \n\t\t Takes a hash as its only argument, and outputs the contents of the object and the\n\t\t node it was retrieved from, if it could be downloaded successfully.\n\texit \n\t\t Terminates the node."
	if !strings.Contains(outputDef, expectedUsage) {
		t.Errorf("Expected blabla in the captured output, but got:\n%s", outputDef)
	}

}

func TestDefault(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 9000)
	node := InitNode(me)
	// Create an input string
	inputString := "notonthelist\n" // Include a newline to simulate pressing Enter

	// Create a buffer to capture the output
	var outputBuffer bytes.Buffer

	// Create an input reader from the input string
	inputReader := strings.NewReader(inputString)

	// Create a CLIConfig with the input reader and the output buffer
	config := CLIConfig{
		Input:  inputReader,
		Output: &outputBuffer,
		Kad:    node,
	}

	// Run the CLI
	go Cli(config)

	time.Sleep(500 * time.Millisecond)

	capturedOutput := outputBuffer.String()
	expectedUsage := "Usage: \n\tput [contents] \n\t\tTakes a single argument, the contents of the file you are uploading, and outputs the\n\t\thash of the object, if it could be uploaded successfully.\n\tget [hash] \n\t\t Takes a hash as its only argument, and outputs the contents of the object and the\n\t\t node it was retrieved from, if it could be downloaded successfully.\n\texit \n\t\t Terminates the node."
	if !strings.Contains(capturedOutput, expectedUsage) {
		t.Errorf("Expected blabla in the captured output, but got:\n%s", capturedOutput)
	}
}

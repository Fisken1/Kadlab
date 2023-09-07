package Kadlab

import (
	d "Kadlab/kademlia"
	"flag"
)

func main() {
	// Define a command-line flag for specifying the node type
	nodeType := flag.String("node-type", "normal", "Specify 'bootstrap' or 'normal' node type")
	bootstrap := flag.String("bootstrapIP", "", "enter ip for bootstrap node")
	flag.Parse()

	var newNode *d.Network

	// Check the value of the nodeType flag
	switch *nodeType {
	case "bootstrap":
		newNode = d.InitBootstrap(7070)
	case "normal":
		newNode, _ = d.InitJoin(*bootstrap, 7070) // Define a function for initializing normal nodes
	default:
		println("Invalid node type specified.")
		return
	}

	// Start listening for incoming connections
	go d.Listen(newNode.Contact.Address, newNode.Contact.Port)

	select {}
}

package kademlia

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Cli(kademlia *Kademlia, port int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">")
	for {
		scanner.Scan()
		text := scanner.Text()

		if len(text) > 0 {
			fmt.Println("we are in if sats")
			input := strings.Fields(text)
			answer := CliHandler(input, kademlia, port)
			fmt.Print(answer + "\n> ")

		} else {
			continue
		}
	}
}

func CliHandler(input []string, node *Kademlia, port int) string {
	fmt.Println("we are in clihandler")
	answer := ""
	switch input[0] {

	case "getContact":
		fmt.Println("getcontact")
		fmt.Println("BUCKETS: ", node.RoutingTable.buckets)

		for _, contacts := range node.RoutingTable.buckets {
			fmt.Println(contacts.list)
		}

	case "put":
		/*
		 (a) put: Takes a single argument, the contents of the file you are uploading, and outputs the
		 hash of the object, if it could be uploaded successfully.
		*/

	case "get":
		/*
		 (b) get: Takes a hash as its only argument, and outputs the contents of the object and the
		 node it was retrieved from, if it could be downloaded successfully.
		*/

	case "exit", "q":
		Terminate()

	default:
		fmt.Println("something in default")
		return "Operation: >>" + input[0] + "<< not found." + "\n" + Usage()
	}
	return answer
}

func Terminate() {
	fmt.Print("Exiting...")
	os.Exit(0)
}

func Usage() string {
	return "Usage: \n\tput [contents] \n\t\tTakes a single argument, the contents of the file you are uploading, and outputs the\n\t\thash of the object, if it could be uploaded successfully.\n\tget [hash] \n\t\t Takes a hash as its only argument, and outputs the contents of the object and the\n\t\t node it was retrieved from, if it could be downloaded successfully.\n\texit \n\t\t Terminates the node."
}

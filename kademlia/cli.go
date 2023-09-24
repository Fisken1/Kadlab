package kademlia

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Cli(kademlia *Kademlia) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("KADEMLIA> ")
	for {
		scanner.Scan()
		text := scanner.Text()

		if len(text) > 0 {
			input := strings.Fields(text)
			answer := CliHandler(input, kademlia)
			fmt.Print(answer + "KADEMLIA> ")

		} else {
			fmt.Print("KADEMLIA> ")
		}
	}
}

func CliHandler(input []string, node *Kademlia) string {
	answer := ""
	switch input[0] {

	//debug
	case "printADDRESS":
		fmt.Println(node.RoutingTable.me.Address)

	//debug
	case "getContact":
		fmt.Println("getcontact")
		fmt.Println("BUCKETS: ", node.RoutingTable.buckets)

		for i, a := range node.RoutingTable.buckets {
			if a.list != nil { // Check if the list is not nil
				for e := a.list.Front(); e != nil; e = e.Next() {
					if e.Value != nil {
						fmt.Println("value in bucket", a, "is", e.Value)
					}
				}
			}
			i++
		}

	case "put":
		inputStrings := input[1:]
		fmt.Println("input strings cli: ", inputStrings)

		var concatenatedString string

		for _, str := range inputStrings {
			concatenatedString += str
		}

		fmt.Println("input strings cli: ", concatenatedString)

		data := []byte(concatenatedString)

		fmt.Println("trying to store..?")
		hash := node.Store(data)

		if hash != "0" {
			answer = hash
		} else {
			answer = "Error..."
		}
		/*
		 (a) put: Takes a single argument, the contents of the file you are uploading, and outputs the
		 hash of the object, if it could be uploaded successfully.
		*/

	case "get":
		hash := input[1]
		_, data := node.LookupData(hash)

		if data != "" {
			answer = data
		} else {
			answer = "Error..."
		}

		/*
		 (b) get: Takes a hash as its only argument, and outputs the contents of the object and the
		 node it was retrieved from, if it could be downloaded successfully.
		*/

	//debug
	case "printID":
		fmt.Println(node.RoutingTable.me.ID.String())

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

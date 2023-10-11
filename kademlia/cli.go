package kademlia

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type CLIConfig struct {
	Input  io.Reader
	Output io.Writer
	Kad    *Kademlia
}

// Cli runs the Kademlia command-line interface.
func Cli(config CLIConfig) {
	scanner := bufio.NewScanner(config.Input)
	fmt.Println("KADEMLIA> ")
	for {
		scanner.Scan()
		text := scanner.Text()

		if len(text) > 0 {
			input := strings.Fields(text)
			answer := CliHandler(input, config.Kad)
			config.Output.Write([]byte(answer))
			fmt.Fprint(config.Output, "KADEMLIA> ")
		} else {
			fmt.Fprint(config.Output, "KADEMLIA> ")
		}
	}
}

func CliHandler(input []string, node *Kademlia) string {
	answer := ""

	switch input[0] {
	case "Uploaded":
		var ToReturn []string
		uploaded := node.storagehandler.getUploadedData()
		for _, data := range uploaded {
			//fmt.Println("Uploaded:", data)
			ToReturn = append(ToReturn, data)
		}
		answer = strings.Join(ToReturn, ", ")

	case "getContact":
		var contactsToReturn []string
		//fmt.Println("getcontact")
		//fmt.Println("BUCKETS: ", node.RoutingTable.buckets)

		for i, a := range node.RoutingTable.buckets {
			if a.list != nil { // Check if the list is not nil
				for e := a.list.Front(); e != nil; e = e.Next() {
					if e.Value != nil {
						//fmt.Println("value in bucket", a, "is", e.Value)
						contactStr := e.Value.(Contact).ID.String()
						contactsToReturn = append(contactsToReturn, contactStr)
					}
				}
			}
			i++
		}
		answer = strings.Join(contactsToReturn, ", ")
		//
		//
	case "put":
		inputStrings := input[1:]

		var concatenatedString string

		for _, str := range inputStrings {
			concatenatedString += str
		}

		data := []byte(concatenatedString)

		hash, _ := node.Store(data)

		if hash != "0" {
			answer = hash + "\n"
			//fmt.Println("Stored")
		} else {
			answer = "Error..." + "\n"
		}

		/*
		 (a) put: Takes a single argument, the contents of the file you are uploading, and outputs the
		 hash of the object, if it could be uploaded successfully.
		*/

	case "get":
		hash := input[1]

		//contact, data := node.LookupData(hash)

		_, contact, data, err := node.LookupData(hash)
		if err != nil {
			answer = "Error..." + "\n"
		} else {
			if data.Value != nil {
				answer = "Found data: " + string(data.Value) + " from contact at addres : " + contact.Address + ". with ID: " + contact.ID.String() + "\n"

			} else {
				answer = "Did not find " + hash + " at any node" + "\n"
			}
		}

		/*
		 (b) get: Takes a hash as its only argument, and outputs the contents of the object and the
		 node it was retrieved from, if it could be downloaded successfully.
		*/
	case "forget":
		hash := input[1]
		b := node.ForgetData(hash)

		if b {
			answer = "Successfully forgot data"
		} else {
			answer = "Failed to forget data"
		}

	case "printID": //debug
		answer = node.RoutingTable.me.ID.String()
		//fmt.Println(node.RoutingTable.me.ID.String())

	case "exit", "q":
		Terminate()

	default:
		//fmt.Println("something in default")
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

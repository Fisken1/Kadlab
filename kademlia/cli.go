package kademlia

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Cli(kademlia *Kademlia, exit chan int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		read, _ := reader.ReadString('\n')
		fmt.Println(CliHandler(strings.Fields(read), kademlia, exit, reader))
	}
}

func CliHandler(s []string, kademlia *Kademlia, exit chan int, reader *bufio.Reader) string {
	if len(s) == 0 {
		return ""
	}
	switch operation := s[0]; operation {
	case "ping":
		//contact := kademlia.LookupContactWithIP(s[0])
		//kademlia.RoutingTable.me.SendPingMessage(contact)
	default:
		return "Operation: " + operation + " not found"
	}
	return "good"
}

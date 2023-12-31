package main

import (
	"Kadlab/kademlia"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	fmt.Println("IP: ", GetOutboundIP())

	node, err := kademlia.InitJoin(GetOutboundIP().String(), 5000)
	if err != nil {
		os.Exit(0)
	}
	config := kademlia.CLIConfig{
		Input:  os.Stdin,
		Output: os.Stdout,
		Kad:    node,
	}
	go kademlia.Cli(config)
	go node.TTLRefresher(30) //Disabled for now to prevent annoying prints
	go node.TTLCleaner(10)   //Disabled for now to prevent annoying prints
	go kademlia.InitHTTPInterface(node, node.RoutingTable)
	select {}

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

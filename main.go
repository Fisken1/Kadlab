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

	go kademlia.Cli(node, 9090)

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

package main

import (
	"Kadlab/kademlia"
	"log"
	"net"
	"strings"
)

func main() {
	node, _ := kademlia.InitJoin(GetOutboundIP().String(), GetBootstrapIP(GetOutboundIP().String()), 8080)
	kademlia.Cli(node, 9090)
	select {}
}

// GetBootstrapIP Check if a node is bootstrap or not, this is hardcoded.
func GetBootstrapIP(ip string) string {
	stringList := strings.Split(ip, ".")
	value := stringList[1]
	bootstrapIP := "130." + value + ".64.25:8080" // some arbitrary IP address hard coded to be bootstrap
	return bootstrapIP
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

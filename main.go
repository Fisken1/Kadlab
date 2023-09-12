package main

import (
	"Kadlab/kademlia"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {

	port := 8080
	ipAddress := GetOutboundIP()
	ipString := ipAddress.String() + ":" + strconv.Itoa(port)
	ipBootstrap := GetBootstrapIP(ipString)

	fmt.Println("your IP is:", ipString)
	fmt.Println("bootstrap IP:", ipBootstrap)
	fmt.Println("port:", port)

	if ipAddress.String() == ipBootstrap {
		bootstrap := kademlia.InitBootstrap(ipBootstrap, port)
		go bootstrap.Listen(*bootstrap.Contact)
	} else {
		node, _ := kademlia.InitJoin(ipAddress.String(), port)
		go node.Listen(*node.Contact)
	}

	select {}

}

// Get preferred outbound ip of this machine
// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// Check if a node is bootstrap or not, this is hardcoded.
func GetBootstrapIP(ip string) string {
	stringList := strings.Split(ip, ".")
	value := stringList[1]
	bootstrapIP := "172." + value + ".0.2:8081" // some arbitrary IP address hard coded to be bootstrap
	return bootstrapIP
}

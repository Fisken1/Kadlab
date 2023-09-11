package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {

	port := "8081"
	ipAddress := GetOutboundIP()
	ipString := ipAddress.String() + ":" + port
	ipBootstrap := GetBootstrapIP(ipString)

	fmt.Println("your IP is:", ipAddress)
	fmt.Println("bootstrap IP:", ipBootstrap)
	fmt.Println("port:", port)

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
	bootstrapIP := "172." + value + ".0.2:10001" // some arbitrary IP address hard coded to be bootstrap
	return bootstrapIP
}

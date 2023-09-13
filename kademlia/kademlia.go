package kademlia

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Hashmap      map[string][]byte
	net          *Network
	alpha        int
}

func InitNode(me Contact) *Kademlia {
	node := &Kademlia{
		RoutingTable: NewRoutingTable(me),
		Hashmap:      make(map[string][]byte),
		alpha:        10,
	}
	node.net = &Network{*node}
	fmt.Print("INITNODE", me.distance, me.ID, me.Address)
	return node
}

func InitJoin(port int) (*Kademlia, error) {
	ipAddress := GetOutboundIP()
	ipString := ipAddress.String() + ":" + strconv.Itoa(port)
	ipBootstrap := GetBootstrapIP(ipString)

	if ipAddress.String() == ipBootstrap {
		bootstrap := InitNode(NewContact(NewRandomKademliaID(), ipBootstrap, port))
		go bootstrap.net.Listen(bootstrap.RoutingTable.me)
		return bootstrap, nil
	} else {
		node := InitNode(NewContact(NewRandomKademliaID(), ipAddress.String(), port))
		node.RoutingTable.me.Port = port
		go Cli(node, make(chan int))
		go node.net.Listen(node.RoutingTable.me)

		return node, nil
	}
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

// GetBootstrapIP Check if a node is bootstrap or not, this is hardcoded.
func GetBootstrapIP(ip string) string {
	stringList := strings.Split(ip, ".")
	value := stringList[1]
	bootstrapIP := "172." + value + ".0.2:8081" // some arbitrary IP address hard coded to be bootstrap
	return bootstrapIP
}

// LookupContact performs a contact lookup for a given key.
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	//TODO
	return nil
}

// performNodeLookup performs a FIND_NODE or FIND_VALUE RPC to a contact.
func (kademlia *Kademlia) performNodeLookup(contact *Contact, target *Contact) []Contact {
	// Implement RPC logic here and return a list of contacts received
	// TODO
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

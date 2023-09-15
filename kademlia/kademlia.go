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
	k            int
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
		bootstrap.RoutingTable.me.Port = port
		go bootstrap.net.Listen(bootstrap.RoutingTable.me)
		go Cli(bootstrap, make(chan int))
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

/*
 * The search begins by selecting alpha contacts from the non-empty k-bucket closest to the bucket appropriate to the key being searched on.
 * If there are fewer than alpha contacts in that bucket, contacts are selected from other buckets. The contact closest to the target key, closestNode, is noted.
 *
 * The first alpha contacts selected are used to create a shortlist for the search.
 *
 * The node then sends parallel, asynchronous FIND_* RPCs to the alpha contacts in the shortlist.
 * Each contact, if it is live, should normally return k triples. If any of the alpha contacts fails to reply, it is removed from the shortlist, at least temporarily.
 *
 * The node then fills the shortlist with contacts from the replies received. These are those closest to the target.
 * From the shortlist it selects another alpha contacts. The only condition for this selection is that they have not already been contacted. Once again a FIND_* RPC is sent to each in parallel.
 *
 * Each such parallel search updates closestNode, the closest node seen so far.
 *
 * The sequence of parallel searches is continued until either no node in the sets returned is closer than the closest node already seen or the initiating node has accumulated k probed and known to be active contacts.
 *
 * If a cycle doesn't find a closer node, if closestNode is unchanged, then the initiating node sends a FIND_* RPC to each of the k closest nodes that it has not already queried.
 *
 * At the end of this process, the node will have accumulated a set of k active contacts or (if the RPC was FIND_VALUE) may have found a data value. Either a set of triples or the value is returned to the caller.
 */

// LookupContact performs a contact lookup for a given key.
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	shortlist := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.alpha)}

	// a list of nodes to know which nodes has been probed already
	probed := ContactCandidates{}

	for {
		updateClosest := false
		numProbed := 0

		for i := 0; i < shortlist.Len() && numProbed < kademlia.k; i++ {

			//kademlia.net.SendFindContactMessage(&shortList.contacts[i])  nånting här?

		}
	}
	return shortlist.contacts
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

//func (kademlia *Kademlia) LookupContactWithIP(ip string) Contact {
//}

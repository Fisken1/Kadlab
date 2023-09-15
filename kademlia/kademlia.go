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

	if ipString == ipBootstrap {
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

	// we need to add code that join needs here

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
	shortlist.Sort()
	contactedNodes := make(map[KademliaID]bool)
	// a list of nodes to know which nodes has been probed already
	probed := ContactCandidates{}

	for {
		updateClosest := false
		numProbed := 0

		for i := 0; i < shortlist.Len() && numProbed < kademlia.k; i++ {

			//kademlia.net.SendFindContactMessage(&shortList.contacts[i])  nånting här?

		}
		para(i)
		i++

		//para
	}
	return shortlist.contacts
}

func (kademlia *Kademlia) para(shortlist ContactCandidates, target *Contact, contactedNodes map[KademliaID]bool, indexMultiplier int) {

	closestNode := shortlist.contacts[0] // Initialize closestNode, you need to define its type based on your code.

	k := kademlia.k // Assuming you have a variable k that defines the desired number of active contacts.
	rpcResults := make(chan []Contact, kademlia.alpha)
	contacts := ContactCandidates{}
	for i := kademlia.alpha * indexMultiplier; i < kademlia.alpha*indexMultiplier+kademlia.alpha; i++ {
		newContacts, err := kademlia.net.SendFindContactMessage(&shortlist.contacts[i], target)
		contactedNodes[*shortlist.contacts[i].ID] = true
		if err != nil {
			fmt.Println(err)
		} else {
			rpcResults <- newContacts
		}
	}

	//rpc res fix

	// Process RPC results.

	for i := 0; i < len(shortlist.contacts); i++ {
		receivedContacts := <-rpcResults
		for _, contact := range receivedContacts {

			if _, alreadyContacted := contactedNodes[*contact.ID]; !alreadyContacted {
				go func(contact Contact) {
					contacts, err := kademlia.net.SendFindContactMessage(&contact, target)
					if err != nil {
						//shortlist remove
					} else {
						rpcResults <- contacts
					}
				}(contact)
				contactedNodes[*contact.ID] = true
			}
		}
		shortlist := ContactCandidates{receivedContacts}

		//look if already contacted

		//shortlist.Sort()
		// Update closestNode if a closer node is found.
		// Add the received contacts to the shortlist.
		// Implement this part based on how you handle results in your system.
	}

	for {
		// Check if you have accumulated k active contacts or no closer nodes than closestNode.
		if shortlist.Len() >= k || closestNode == *target {
			break
		}

		// Send FIND_* RPCs to alpha contacts in the shortlist.
		for _, contact := range shortlist.contacts {
			if _, alreadyContacted := contactedNodes[*contact.ID]; !alreadyContacted {
				go func(contact Contact) {
					contacts, err := kademlia.net.SendFindContactMessage(&contact, target)
					if err != nil {
						//shortlist remove
					} else {
						rpcResults <- contacts
					}
				}(contact)
				contactedNodes[*contact.ID] = true
			}
		}
	}

	//receivedContacts := <-rpcResults

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

func (kademlia *Kademlia) LookupNode2(target *Contact) (*Contact, error) {
	// Initialize variables
	k := kademlia.k
	queriedContacts := []Contact{}
	contactedMap := make(map[string]bool)
	closestNode := kademlia.getClosestNode(*target.ID, queriedContacts)
	closestNodeSeen := closestNode
	iterationsWithoutCloserNode := 0

	for {
		// Check if you have accumulated k active contacts or found the closest node.
		if len(queriedContacts) >= k || closestNode.ID == target.ID {
			break
		}

		// Send FIND_NODE RPCs to alpha contacts in the shortlist.
		alphaContacts := kademlia.getAlphaContacts(closestNode, queriedContacts, k, contactedMap)
		foundContacts, err := kademlia.QueryAlphaContacts(alphaContacts, target)

		if err != nil {
			// Handle errors if the RPC fails.
		}

		// Update queriedContacts with the results.
		queriedContacts = append(queriedContacts, foundContacts...)

		// Find the closest node among the newly queried contacts.
		newClosestNode := kademlia.getClosestNode(*target.ID, queriedContacts)

		// Check if the closest node has not changed.
		if newClosestNode.ID == closestNode.ID {
			// If no closer node is found in this cycle, stop the search.
			break
		}

		closestNode = newClosestNode
	}

	return closestNode, nil
}

func (kademlia *Kademlia) getAlphaContacts(node *Contact, queriedContacts []Contact, alpha int, contactedMap map[string]bool) []Contact {
	var alphaContacts []Contact

	// Find the closest contacts to the current node using FindClosestContacts.
	closestContacts := kademlia.RoutingTable.FindClosestContacts(node.ID, alpha)

	// Iterate through the closest contacts and filter out those that have already been contacted.
	for _, neighbor := range closestContacts {
		// Check if the neighbor is not already in queriedContacts.
		if _, alreadyContacted := contactedMap[neighbor.ID.String()]; !alreadyContacted {
			alphaContacts = append(alphaContacts, neighbor)
		}

		// Stop if we have collected enough alpha contacts.
		if len(alphaContacts) >= alpha {
			return alphaContacts
		}
	}

	return alphaContacts
}

func (kademlia *Kademlia) getClosestNode(targetID KademliaID, queriedContacts []Contact) *Contact {
	var closest *Contact
	var minDistance *KademliaID

	// Iterate through the queriedContacts to find the closest node.
	for _, contact := range queriedContacts {
		// Calculate the XOR distance between targetID and the contact's ID using CalcDistance method.
		contact.CalcDistance(&targetID)

		// If closest is nil or the current contact is closer, update closest and minDistance.
		if closest == nil || contact.distance.Less(minDistance) {
			closest = &contact
			minDistance = contact.distance
		}
	}

	return closest
}

func (kademlia *Kademlia) QueryAlphaContacts(alphaContacts []Contact, target *Contact) ([]Contact, error) {
	// Create channels to receive results and errors.
	results := make(chan []Contact, len(alphaContacts))
	errors := make(chan error, len(alphaContacts))

	// Iterate through alphaContacts and send FIND_NODE RPCs in parallel.
	for _, contact := range alphaContacts {
		go func(contact Contact) {
			// Send a FIND_NODE RPC to the contact.
			foundContacts, err := kademlia.net.SendFindContactMessage(&contact, target)
			if err != nil {
				// Handle the error and send it to the errors channel.
				errors <- err
			} else {
				// Send the found contacts to the results channel.
				results <- foundContacts
			}
		}(contact)
	}

	// Collect the results from the channels.
	var foundContacts []Contact
	for i := 0; i < len(alphaContacts); i++ {
		select {
		case contacts := <-results:
			// Add the found contacts to the result slice.
			foundContacts = append(foundContacts, contacts...)
		case err := <-errors:
			// Handle errors here if needed.
			fmt.Printf("Error: %v\n", err)
		}
	}

	return foundContacts, nil
}

package kademlia

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Kademlia struct {
	RoutingTable     *RoutingTable
	Hashmap          map[string][]byte
	net              *Network
	alpha            int
	k                int
	bootstrap        bool
	bootstrapContact *Contact
}

const (
	BootstrapKademliaID = "FFFFFFFF00000000000000000000000000000000"
)

func InitNode(me Contact) *Kademlia {
	node := &Kademlia{
		RoutingTable: NewRoutingTable(me),
		Hashmap:      make(map[string][]byte),
		alpha:        3,
	}
	node.net = &Network{nil}
	fmt.Print("INITNODE", me.ID.String(), me.Address)
	return node
}

func InitJoin(ip string, port int) (*Kademlia, error) {

	if ip+":"+strconv.Itoa(port) == GetBootstrapIP(ip)+":"+strconv.Itoa(port) {
		bootstrap := InitNode(NewContact(NewKademliaID(BootstrapKademliaID), GetBootstrapIP(ip), port))
		bootstrap.bootstrap = true
		bootstrap.RoutingTable.me.Port = port
		go bootstrap.net.Listen(bootstrap.RoutingTable.me)
		return bootstrap, nil
	} else {
		node := InitNode(NewContact(NewRandomKademliaID(), ip, port))
		node.bootstrap = false
		node.RoutingTable.me.Port = port

		contact := NewContact(
			NewKademliaID(BootstrapKademliaID),
			GetBootstrapIP(ip),
			port,
		)
		node.bootstrapContact = &contact
		go node.net.Listen(node.RoutingTable.me)

		return node, nil
	}

	// we need to add code that join needs here
}

// GetBootstrapIP Check if a node is bootstrap or not, this is hardcoded.
func GetBootstrapIP(ip string) string {
	stringList := strings.Split(ip, ".")
	value := stringList[1]
	bootstrapIP := "192." + value + ".1.26" // some arbitrary IP address hard coded to be bootstrap
	return bootstrapIP
}

func (kademlia *Kademlia) fixNetwork() {

	if kademlia.bootstrap {
		fmt.Println("You are the bootstrap node!")
		return

	}

	err := kademlia.net.SendPingMessage(&kademlia.RoutingTable.me, kademlia.bootstrapContact)
	if err != nil {
		return
	}

	kademlia.RoutingTable.AddContact(*kademlia.bootstrapContact)

	contacts, err := kademlia.LookupNode(&kademlia.RoutingTable.me)
	if err != nil {
		return
	}
	for _, contact := range contacts {
		kademlia.RoutingTable.AddContact(contact)
	}

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

func (kademlia *Kademlia) LookupNode(target *Contact) ([]Contact, error) {
	// Initialize variables
	k := kademlia.k
	queriedContacts := []Contact{}
	contactedMap := make(map[string]bool)
	closestNode := kademlia.getClosestNode(*target.ID, queriedContacts)

	for {
		// Check if you have accumulated k active contacts or found the closest node.
		if len(queriedContacts) >= k || closestNode.ID == target.ID {
			break
		}

		// Send FIND_NODE RPCs to alpha contacts in the shortlist.
		alphaContacts := kademlia.getAlphaContacts(closestNode, queriedContacts, k, contactedMap)
		shortlist, err := kademlia.QueryAlphaContacts(alphaContacts, target)

		if err != nil {
			// Handle errors if the RPC fails.
		}

		// Update queriedContacts with the results.
		queriedContacts = append(queriedContacts, shortlist...)

		// Find the closest node among the newly queried contacts.
		newClosestNode := kademlia.getClosestNode(*target.ID, queriedContacts)

		// Check if the closest node has not changed.
		if newClosestNode.ID == closestNode.ID {
			// If no closer node is found in this cycle, stop the search.
			break
		}

		closestNode = newClosestNode
	}
	queriedContacts = kademlia.getKNodes(queriedContacts, target)
	return queriedContacts, nil
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
			if closest != nil {
				fmt.Println("Switched closest contact from : " + closest.ID.String() + " to " + contact.ID.String() + " with a distance of " + contact.distance.String())
			} else {
				fmt.Println("FIRST ENTRY :" + contact.ID.String())
			}

			// Create a separate variable to store the current contact.
			currentContact := contact
			closest = &currentContact
			minDistance = currentContact.distance
			fmt.Println("NOW THE MINDISTANCE IS: " + minDistance.String())
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
			//HANDLE SEEN NODES!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
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

func (kademlia *Kademlia) getKNodes(contacts []Contact, target *Contact) []Contact {
	// Ensure k is within the bounds of the contacts slice.
	k := kademlia.k
	if k <= 0 || k >= len(contacts) {
		return contacts // Return all contacts if k is out of bounds.
	}

	// Calculate XOR distances for all contacts.
	for i := range contacts {
		contacts[i].CalcDistance(target.ID)
	}

	// Sort the contacts based on their XOR distances.
	sort.Slice(contacts, func(i, j int) bool {
		return contacts[i].distance.Less(contacts[j].distance)
	})

	// Return the first k contacts, which are the closest ones.
	return contacts[:k]
}

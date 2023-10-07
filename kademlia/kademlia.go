package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
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
		k:            4,
	}
	node.net = &Network{node: node}
	//fmt.Print("THIS IS HOW THE NODE LOOKS LIKE AFTER WE CREATE IT ", me.ID.String(), me.Address, me.ID)
	return node
}

func InitJoin(ip string, port int) (*Kademlia, error) {
	fmt.Println(GetBootstrapIP(ip))
	fmt.Println(ip + "." + strconv.Itoa(port))
	if ip+":"+strconv.Itoa(port) == GetBootstrapIP(ip)+":"+strconv.Itoa(port) {
		fmt.Println("WE ARE BOOTSTRAP")
		bootstrap := InitNode(NewContact(NewKademliaID(BootstrapKademliaID), GetBootstrapIP(ip), port))
		bootstrap.bootstrap = true
		bootstrap.RoutingTable.me.Port = port
		go bootstrap.net.Listen(bootstrap.RoutingTable.me)

		// Use a channel for synchronization.
		done := make(chan bool)

		// Start a goroutine to close the network and signal when done.
		go func() {
			// Sleep for a while to allow the network to start.
			time.Sleep(1 * time.Second)

			done <- true
		}()

		<-done

		bootstrap.fixNetwork()

		return bootstrap, nil
	} else {
		node := InitNode(NewContact(NewRandomKademliaID(), ip, port))
		node.bootstrap = false
		node.RoutingTable.me.Port = port

		contact := NewContact(
			NewKademliaID(BootstrapKademliaID),
			GetBootstrapIP(ip),
			GetBootstrapPort(),
		)
		node.bootstrapContact = &contact
		go node.net.Listen(node.RoutingTable.me)

		// Use a channel for synchronization.
		done := make(chan bool)

		// Start a goroutine to close the network and signal when done.
		go func() {
			// Sleep for a while to allow the network to start.
			time.Sleep(1 * time.Second)

			done <- true
		}()

		<-done

		node.fixNetwork()

		return node, nil
	}

	// we need to add code that join needs here
}

// GetBootstrapIP Check if a node is bootstrap or not, this is hardcoded.
func GetBootstrapIP(ip string) string {

	/*
		stringList := strings.Split(ip, ".")
		value := stringList[1]
		bootstrapIP := "172." + value + ".0.2" // some arbitrary IP address hard coded to be bootstrap
		return bootstrapIP

	*/
	bootstrapIPForTests := "192.168.1.26" // some arbitrary IP address hard coded to be bootstrap
	return bootstrapIPForTests
}

func GetBootstrapPort() int {
	return 5000
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

	//contacts, err := kademlia.LookupNode(&kademlia.RoutingTable.me)
	contacts, err := kademlia.LookupNode(&kademlia.RoutingTable.me)
	if err != nil {
		return
	}

	fmt.Println("THIS IS US", kademlia.RoutingTable.me.Address+":"+strconv.Itoa(kademlia.RoutingTable.me.Port), "AND THESE ARE THE CONTACTS THAT WE ARE ADDING TO THE ")
	for _, contact := range contacts {
		kademlia.RoutingTable.AddContact(contact)
		fmt.Println("\t\tCONTACT", contact.Address+":"+strconv.Itoa(contact.Port))
	}
}

func (kademlia *Kademlia) Store(data []byte) string {
	type resultStruct struct {
		msg string
		err error
	}

	keyString := hex.EncodeToString(sha1.New().Sum(data))
	target := NewContact(NewKademliaID(keyString), "000.000.00.0", 0000)

	contacts, err := kademlia.LookupNode(&target)

	resultsChannel := make(chan resultStruct, int(math.Min(float64((kademlia.alpha)), float64(len(contacts)))))
	if err != nil {
		panic(err)
	}
	//TODO choose duplicate amount currently alpha duplicates
	for i := 0; i < len(contacts) && i < kademlia.alpha; i++ {
		if contacts[i].ID.String() == kademlia.RoutingTable.me.ID.String() {

			kademlia.handleStoreMessage(target.ID.String(), data)
			fmt.Println("Stored at Self")
			resultsChannel <- resultStruct{"Stored at self", nil}
		} else {
			msg, err := kademlia.net.SendStoreMessage(&kademlia.RoutingTable.me, &contacts[i], target.ID.String(), data)
			if err != nil {
				fmt.Println("error during storing")
				// Handle the error and send it to the errors channel.
				resultsChannel <- resultStruct{"", err}
			} else {
				resultsChannel <- resultStruct{msg, nil}
			}
		}

	}
	for i := 0; i < int(math.Min(float64((kademlia.alpha)), float64(len(contacts)))); i++ {
		select {

		case results := <-resultsChannel:
			if results.err != nil {
				fmt.Println("Store response message: ", results.msg)
			} else {
				fmt.Println(err)
			}

		}
	}

	return target.ID.String()

}

// TODO not currently used
func (kademlia *Kademlia) getAlphaContacts(node *Contact, alpha int, contactedMap map[string]bool) []Contact {
	var alphaContacts []Contact

	// Find the closest contacts to the current node using FindClosestContacts.
	closestContacts := kademlia.RoutingTable.FindClosestContacts(node.ID, alpha)
	fmt.Println("now we are in getAlphaContacts and the amount of contacts found are: ", closestContacts, "and these are the contacts in the map", contactedMap)
	// Iterate through the closest contacts and filter out those that have already been contacted.
	for _, neighbor := range closestContacts {
		// Check if the neighbor is not already in queriedContacts.
		if _, alreadyContacted := contactedMap[neighbor.ID.String()]; !alreadyContacted {
			alphaContacts = append(alphaContacts, neighbor)
			fmt.Println("Alphacontact add:", neighbor.Address+":"+strconv.Itoa(neighbor.Port))
		}

		// Stop if we have collected enough alpha contacts.
		if len(alphaContacts) >= alpha {
			return alphaContacts
		}
	}

	return alphaContacts
}
func (kademlia *Kademlia) getAlphaContactsMap(contacts map[string]Contact) (map[string]Contact, map[string]Contact) {
	alphaContacts := make(map[string]Contact)
	remainingContacts := make(map[string]Contact)

	i := 0
	for _, contact := range contacts {
		// Check if the neighbor is not already in queriedContacts.
		if i < kademlia.alpha {
			alphaContacts[contact.ID.String()] = contact
			i++
		} else {
			remainingContacts[contact.ID.String()] = contact
		}

	}

	return alphaContacts, remainingContacts
}

func (kademlia *Kademlia) getClosestNode(targetID KademliaID, contacts []Contact) *Contact {
	var closest *Contact
	var minDistance *KademliaID

	//fmt.Println("amount of contacts", len(contacts))

	// Iterate through the queriedContacts to find the closest node.
	for _, contact := range contacts {
		// Calculate the XOR distance between targetID and the contact's ID using CalcDistance method.
		contact.CalcDistance(&targetID)

		//fmt.Print("WOW THE CONTACT WE ARE LOOKING AT JUST GOT A DISTANCE", contact.ID, contact.distance)

		// If closest is nil or the current contact is closer, update closest and minDistance.
		if closest == nil || contact.distance.Less(minDistance) {
			if closest != nil {
				//fmt.Println("Switched closest contact from : " + closest.ID.String() + " to " + contact.ID.String() + " with a distance of " + contact.distance.String())
			} else {
				//fmt.Println("FIRST ENTRY :" + contact.ID.String())
			}

			// Create a separate variable to store the current contact.
			currentContact := contact
			closest = &currentContact
			minDistance = currentContact.distance
			//fmt.Println("NOW THE MINDISTANCE IS: " + minDistance.String())
		}
	}

	return closest
}

func (kademlia *Kademlia) LookupNode(target *Contact) ([]Contact, error) {
	// Initialize variables
	var lU lookUpper

	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.alpha)

	lU.initLookUp(*kademlia, *target, closestContacts, "FIND_NODE")
	finalResult, _, _ := lU.startLookup()
	return finalResult, nil

}
func (kademlia *Kademlia) LookupData(hash string) ([]Contact, Contact, []byte, error) {
	//Exist at self
	value, exists := kademlia.Hashmap[hash]
	if exists {
		return []Contact{kademlia.RoutingTable.me}, kademlia.RoutingTable.me, value, nil
	}
	// Initialize variables
	var lU lookUpper
	target := NewContact(NewKademliaID(hash), "000.000.00.00", 0000)
	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.alpha)
	lU.initLookUp(*kademlia, target, closestContacts, "FIND_DATA")
	finalResult, con, val := lU.startLookup()

	return finalResult, con, val, nil

}

func (kademlia *Kademlia) QueryContacts(contacts []Contact, target *Contact) ([]Contact, error) {

	type resultStruct struct {
		contacts []Contact
		err      error
	}

	foundContacts := make(map[string]Contact)

	//Create channels to receive results and errors.
	resultsChannel := make(chan resultStruct, len(contacts))

	for _, contact := range contacts {
		//fmt.Println("this is the contact we are looking at now", contact)
		go func(contact Contact) {

			// Send a FIND_NODE RPC to the contact.

			newContacts, err := kademlia.net.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target)

			if err != nil {
				fmt.Println("error during querry contact. ", err)
				resultsChannel <- resultStruct{nil, err}
			} else {

				kademlia.RoutingTable.AddContact(contact)
				resultsChannel <- resultStruct{newContacts, nil}
			}

		}(contact)
	}

	// Collect the results from the channels.
	for i := 0; i < len(contacts); i++ {
		select {
		case results := <-resultsChannel:

			if results.err != nil {

			} else {
				for _, c := range results.contacts {

					foundContacts[c.ID.String()] = c

				}
			}

		}
	}

	//Map removes duplicates
	contactList := getListFromMap(foundContacts)
	return contactList, nil
}
func (kademlia *Kademlia) QueryContactsForValue(contacts []Contact, hash string) ([]Contact, Contact, []byte, error) {
	type resultStruct struct {
		contacts []Contact
		value    []byte
		err      error
	}

	var foundValue []byte
	var keeperContact Contact
	foundContacts := make(map[string]Contact)

	//Create channels to receive results and errors.
	resultsChannel := make(chan resultStruct, len(contacts))

	// Iterate through alphaContacts and send FIND_VALUE RPCs in parallel.
	for _, contact := range contacts {
		go func(contact Contact) {
			// Send a FIND_VALUE RPC to the contact.

			newContacts, value, err := kademlia.net.SendFindDataMessage(&kademlia.RoutingTable.me, &contact, hash)

			if err != nil {

				fmt.Printf("Error: %v\n", err)
				resultsChannel <- resultStruct{nil, nil, err}

			} else {
				resultsChannel <- resultStruct{newContacts, value, err}

				//If value is found
				if value != nil {

					keeperContact = newContacts[0]
					foundValue = value
				}

			}

		}(contact)
	}

	// Collect the results from the channels.
	for i := 0; i < len(contacts); i++ {
		select {
		case results := <-resultsChannel:

			if results.err != nil {

			} else {
				for _, c := range results.contacts {

					foundContacts[c.ID.String()] = c

				}
			}

		}
	}
	//Map removes duplicates
	contactList := getListFromMap(foundContacts)
	return contactList, keeperContact, foundValue, nil
}

// TODO not currently used
// Trim list so that the length is k
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
func (kademlia *Kademlia) handleStoreMessage(keyString string, data []byte) {
	fmt.Println("Storing")
	kademlia.Hashmap[keyString] = data

}

func getListFromMap(contactMap map[string]Contact) []Contact {
	var contacts []Contact
	for _, contact := range contactMap {
		contacts = append(contacts, contact)
	}
	return contacts
}

type nodeData struct {
	contact Contact
	//distance KademliaID
	probed bool
	value  []byte
}
type lookUpper struct {
	nodeList      []nodeData
	target        Contact
	kademlia      Kademlia
	mode          string
	storedContact Contact
	foundValue    []byte
}

func (lU *lookUpper) initLookUp(kademlia Kademlia, target Contact, closestNodes []Contact, mode string) {
	lU.kademlia = kademlia
	lU.target = target
	lU.mode = mode
	var nData nodeData
	fmt.Println("starting lookup with target: ", lU.target.ID.String())
	//cList := lU.distanceSort(closestNodes, &target)
	var nList = make([]nodeData, len(closestNodes))
	for i := 0; i < len(closestNodes); i++ {
		nData = nodeData{closestNodes[i], false, nil}
		nData.contact.CalcDistance(target.ID)
		nList[i] = nData

	}
	sort.Slice(nList, func(i, j int) bool {
		return nList[i].contact.distance.Less(nList[j].contact.distance)
	})
	lU.nodeList = nList

}
func (lU *lookUpper) startLookup() ([]Contact, Contact, []byte) {

	var alphaCont []Contact
	var newContacts []Contact

	var completeBool bool
	var valueContact Contact
	var foundValue []byte
	for {
		completeBool = lU.completeCheck()
		if completeBool {
			break
		}
		alphaCont = lU.getAlphaForContact(lU.kademlia.alpha)
		if lU.mode == "FIND_NODE" {
			newContacts, _ = lU.kademlia.QueryContacts(alphaCont, &lU.target)
		}
		if lU.mode == "FIND_DATA" {
			newContacts, valueContact, foundValue, _ = lU.kademlia.QueryContactsForValue(alphaCont, lU.target.ID.String())
			if foundValue != nil {
				lU.storedContact = valueContact
				lU.foundValue = foundValue
			}
		}
		for _, contact := range newContacts {
			lU.insert(contact)
		}
	}
	finalResult := lU.getNList(lU.kademlia.k)
	return finalResult, lU.storedContact, lU.foundValue
}
func (lU *lookUpper) insert(contact Contact) {
	var nList []nodeData = lU.nodeList
	var match = false
	for _, nodedata := range lU.nodeList {
		if nodedata.contact.ID.String() == contact.ID.String() {
			match = true
		}
	}
	if !match {
		c := nodeData{contact, false, nil}
		c.contact.CalcDistance(lU.target.ID)
		nList = append(nList, c)
		sort.Slice(nList, func(i, j int) bool {
			return nList[i].contact.distance.Less(nList[j].contact.distance)
		})
		lU.nodeList = nList
	}

}
func (lU *lookUpper) completeCheck() bool {
	var cList []Contact
	if lU.foundValue != nil && lU.mode == "FIND_DATA" {
		return true
	}
	for i := 0; i < lU.kademlia.k && i < len(lU.nodeList); i++ {
		if lU.nodeList[i].probed == false {
			return false
		}
		cList = append(cList, lU.nodeList[i].contact)

	}
	return true

}
func (lU *lookUpper) getAlphaForContact(amount int) []Contact {
	var contacts []Contact
	n := 0
	for i := 0; i < lU.kademlia.k && i < len(lU.nodeList); i++ {
		if lU.nodeList[i].probed == false {
			if n < amount {
				lU.nodeList[i].probed = true
				contacts = append(contacts, lU.nodeList[i].contact)
				n++
			}
		}

	}
	return contacts

}

func (lU *lookUpper) getList() []Contact {
	var contacts []Contact
	for _, node := range lU.nodeList {
		contacts = append(contacts, node.contact)
	}
	return contacts
}
func (lU *lookUpper) getNList(length int) []Contact {
	var contacts []Contact
	for i := 0; i < length && i < len(lU.nodeList); i++ {

		contacts = append(contacts, lU.nodeList[i].contact)
	}
	return contacts
}

func (lU *lookUpper) printList() {
	for _, node := range lU.nodeList {
		fmt.Println(node.contact.ID)
	}
}

// TODO remove, not currently used
func (lU *lookUpper) distanceSort(contacts []Contact, target *Contact) []Contact {
	// Ensure k is within the bounds of the contacts slice.

	// Calculate XOR distances for all contacts.
	for i := range contacts {
		contacts[i].CalcDistance(target.ID)

	}

	// Sort the contacts based on their XOR distances.
	sort.Slice(contacts, func(i, j int) bool {
		return contacts[i].distance.Less(contacts[j].distance)
	})

	// Return the first k contacts, which are the closest ones.
	return contacts
}

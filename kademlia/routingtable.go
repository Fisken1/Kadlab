package kademlia

import (
	"fmt"
	"sync"
)

const bucketSize = 20

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket
	mutex   sync.RWMutex
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

// AddContact add a new contact to the correct Bucket
// Dont use this directly, use kademlia.AddContact(contact Contact) instead.
func (routingTable *RoutingTable) AddContact(contact Contact) (bool, Contact) {
	addedToFront := true
	bucketIndex := routingTable.getBucketIndex(contact.ID)

	bucket := routingTable.buckets[bucketIndex]
	fmt.Printf(" %-10s  %-10s  %-10s %-10s %-10s %-10x \n", "Attempting to add ", contact.Address, " to this ", routingTable.me.Address, " at bucketIndex: ", bucketIndex)
	bucket.AddContact(contact)
	frontContact := bucket.list.Front().Value.(Contact)
	//If the contact is not at the front then the first should be contacted
	if bucket.list.Front().Value.(Contact).ID.String() != contact.ID.String() {
		//fmt.Println()
		addedToFront = false
		fmt.Printf(" %-10s  %-10s  %-10s %-10s %-10s %-10x %-10s %-10x\n", "Adding failed for ", contact.Address, " to this ", routingTable.me.Address, " at bucketIndex: ", bucketIndex, " Reason: bucket is full. len(bucket): ", bucket.list.Len())

	} else {
		fmt.Printf(" %-10s  %-10s  %-10s %-10s %-10s %-10x %-10s %-10x \n", "Successfully added ", contact.Address, " to this ", routingTable.me.Address, " at bucketIndex: ", bucketIndex, " Length of bucket: ", bucket.list.Len())
	}
	return addedToFront, frontContact

}

func (routingTable *RoutingTable) removeContactAtFront(contact Contact) bool {
	bucketIndex := routingTable.getBucketIndex(contact.ID)

	bucket := routingTable.buckets[bucketIndex]
	bucket.removeFrontElement()
	return true
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	routingTable.lock()
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {

		if bucketIndex-i >= 0 {

			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {

			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}
	routingTable.unLock()

	return candidates.GetContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}

func (routingTable *RoutingTable) lock() {
	routingTable.mutex.Lock()
}

func (routingTable *RoutingTable) unLock() {
	routingTable.mutex.Unlock()
}

package kademlia

import (
	"fmt"
	"sync"
	"time"
)

type StorageHandler struct {
	valueMap   map[string]StorageData
	mutex      sync.RWMutex
	defaultTTL time.Duration
}
type StorageData struct {
	Key        string    `json:"Key"`
	Value      []byte    `json:"Data,omitempty"`
	TimeToLive time.Time `json:"TTL,omitempty"`
	Original   bool      `json:"Original,omitempty"`
}

func (storagehandler *StorageHandler) initStorageHandler() {
	storagehandler.valueMap = make(map[string]StorageData)
	storagehandler.defaultTTL = time.Duration(time.Second * 60) //Change if requierd
}
func (storagehandler *StorageHandler) lock() {
	storagehandler.mutex.Lock()

}

func (storagehandler *StorageHandler) unLock() {
	storagehandler.mutex.Unlock()
}

func (storagehandler *StorageHandler) store(storageData StorageData) {

	storagehandler.lock()
	storagehandler.valueMap[storageData.Key] = storageData
	storagehandler.unLock()
}

//returns value if it exist and a boolean depending if the key got a match.
func (storagehandler *StorageHandler) getValue(key string) (StorageData, bool) {
	storagehandler.lock()
	valuestruct, exists := storagehandler.valueMap[key]
	if exists && valuestruct.Original {
		valuestruct.TimeToLive = time.Now().Add(storagehandler.defaultTTL) //Updates TTL
		storagehandler.valueMap[valuestruct.Key] = valuestruct
	}
	storagehandler.unLock()
	return valuestruct, exists
}
func (storagehandler *StorageHandler) forgetData(key string) bool {
	storagehandler.lock()

	valuestruct, b := storagehandler.valueMap[key]
	if b {
		fmt.Println("forgetting. key: ", valuestruct.Key, " value: ", valuestruct.Value)

		valuestruct.Original = false
		storagehandler.valueMap[valuestruct.Key] = valuestruct
	}

	storagehandler.unLock()
	return b
}

func (storagehandler *StorageHandler) getOriginalData() []StorageData {
	originalData := []StorageData{}
	storagehandler.lock()
	for _, valueStruct := range storagehandler.valueMap {

		if valueStruct.Original {
			originalData = append(originalData, valueStruct)
		}

	}
	storagehandler.unLock()
	return originalData

}

func (storagehandler *StorageHandler) clearExpiredData() {

	storagehandler.lock()
	for _, valueStruct := range storagehandler.valueMap {
		if valueStruct.TimeToLive.Before(time.Now()) { //Consider adding some extra time to give some leeway
			fmt.Println("Removing value with past TTL, key: ", valueStruct.Key, " value: ", valueStruct.Value)
			delete(storagehandler.valueMap, valueStruct.Key)
		}
	}
	storagehandler.unLock()

}

func (storagehandler *StorageHandler) printStoredData() {
	storagehandler.lock()
	for _, data := range storagehandler.valueMap {
		fmt.Printf(" %-10s  %-10s  %-10s %-10s %-10s %-10s %-10s %-10t \n", "Key: ", data.Key, " Value: ", data.Value, " TTL: ", data.TimeToLive, " Original: ", data.Original)
	}
	storagehandler.unLock()
}

package kademlia

import (
	"sync"
	"time"
)

type StorageHandler struct {
	valueMap     map[string]StorageData
	mutex        sync.RWMutex
	defaultTTL   time.Duration
	uploadedData map[string]string
}
type StorageData struct {
	Key        string    `json:"Key"`
	Value      []byte    `json:"Data,omitempty"`
	TimeToLive time.Time `json:"TTL,omitempty"`
}

func (storagehandler *StorageHandler) initStorageHandler() {
	storagehandler.valueMap = make(map[string]StorageData)
	storagehandler.uploadedData = make(map[string]string)
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

// returns value if it exist and a boolean depending if the key got a match.
func (storagehandler *StorageHandler) getValue(key string) (StorageData, bool) {
	storagehandler.lock()
	valuestruct, exists := storagehandler.valueMap[key]
	if exists {
		valuestruct.TimeToLive = time.Now().Add(storagehandler.defaultTTL) //Updates TTL
		storagehandler.valueMap[valuestruct.Key] = valuestruct
	}
	storagehandler.unLock()
	return valuestruct, exists
}
func (storagehandler *StorageHandler) forgetData(key string) bool {
	storagehandler.lock()
	_, exists := storagehandler.uploadedData[key]
	if exists {
		delete(storagehandler.uploadedData, key)
	}

	storagehandler.unLock()
	return exists
}

func (storagehandler *StorageHandler) getUploadedData() []string {
	storagehandler.lock()
	uploads := []string{}

	for _, data := range storagehandler.uploadedData {
		uploads = append(uploads, data)
	}

	storagehandler.unLock()
	return uploads
}

func (storagehandler *StorageHandler) setTTL(storageData StorageData) bool {
	storagehandler.lock()
	valuestruct, exists := storagehandler.valueMap[storageData.Key]
	if exists {
		valuestruct.TimeToLive = storageData.TimeToLive
		storagehandler.valueMap[valuestruct.Key] = valuestruct
	}

	storagehandler.unLock()
	return exists
}

func (storagehandler *StorageHandler) clearExpiredData() []StorageData {

	storagehandler.lock()
	oldData := []StorageData{}
	for _, valueStruct := range storagehandler.valueMap {
		if valueStruct.TimeToLive.Before(time.Now()) { //Consider adding some extra time to give some leeway
			oldData = append(oldData, valueStruct)
			delete(storagehandler.valueMap, valueStruct.Key)
		}
	}
	storagehandler.unLock()
	return oldData

}

func (storagehandler *StorageHandler) setUploader(key string) {
	storagehandler.lock()
	storagehandler.uploadedData[key] = key
	storagehandler.unLock()
}
func (storagehandler *StorageHandler) getStoredData() []StorageData {
	storedData := []StorageData{}
	storagehandler.lock()
	for _, data := range storagehandler.valueMap {
		storedData = append(storedData, data)
	}
	storagehandler.unLock()
	return storedData
}

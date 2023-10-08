package kademlia

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type Network struct {
	node   *Kademlia
	server *net.UDPConn
}

// KademliaMessage represents a Kademlia message.
type KademliaMessage struct {
	Type     string       `json:"Type"`
	Sender   *Contact     `json:"SenderID"`
	Receiver *Contact     `json:"NodeID"`
	Target   *Contact     `json:"TargetID"`
	Data     *StorageData `json:"Data,omitempty"`
	Contacts []Contact    `json:"Contacts,omitempty"`
}

func CreateKademliaMessage(messageType string, sender, receiver, target *Contact, data *StorageData, contacts []Contact) KademliaMessage {
	return KademliaMessage{
		Type:     messageType,
		Sender:   sender,
		Receiver: receiver,
		Target:   target,
		Data:     data,
		Contacts: contacts,
	}
}

func (network *Network) Listen(contact Contact) {

	addr := contact.Address + ":" + strconv.Itoa(contact.Port)

	listenAdrs, _ := net.ResolveUDPAddr("udp", addr)
	servr, err := net.ListenUDP("udp", listenAdrs)
	network.server = servr
	if err != nil {
		fmt.Println("BIG ERROR! This is the addr it tried to listen to: "+listenAdrs.String(), err)
		defer servr.Close()
	}
	fmt.Println("Listening on: " + listenAdrs.String() + ". Id: " + contact.ID.String() + "\n\n")

	for {
		buf := make([]byte, 65536)
		rlen, rem, err := servr.ReadFromUDP(buf)
		if err != nil {
			//fmt.Println("Error reading from UDP:", err)
			continue
		}

		// Process the received data by passing it to the Dispatcher function.
		go func(connection *net.UDPConn) {
			response, err := network.Dispatcher(buf[:rlen])
			if err != nil {
				log.Printf("Failed to handle response message: %v\n", err)
				return
			}
			connection.WriteToUDP(response, rem)

		}(servr)

	}
}

// Dispatcher is responsible for routing incoming messages to their respective handlers.
func (network *Network) Dispatcher(data []byte) ([]byte, error) {
	var msg KademliaMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("In Dispatcher: Error decoding message:", err)
		return nil, err
	}

	fmt.Printf(" %-10s  %-10s  %-10s %-10s %-10s %-10s \n", "This node: ", msg.Receiver.Address, " Got a message from: ", msg.Sender.Address, " Of the type: ", msg.Type)

	switch msg.Type {
	case "PING":
		network.node.RoutingTable.AddContact(*msg.Sender)
		response := CreateKademliaMessage(
			"PONG",
			msg.Receiver,
			msg.Sender,
			nil,
			nil,
			[]Contact{},
		)
		fmt.Println("Sending: ", response.Type, " to: ", msg.Sender)
		msgToSend, err := json.Marshal(response)

		return msgToSend, err

	case "FIND_NODE":
		network.node.RoutingTable.AddContact(*msg.Sender)
		closestNodes := network.node.RoutingTable.FindClosestContacts(msg.Receiver.ID, network.node.k)
		fmt.Println("closestNodes: ", closestNodes)
		response := CreateKademliaMessage(
			"FIND_NODE_RESPONSE",
			msg.Receiver,
			msg.Sender,
			msg.Target,
			nil,
			closestNodes,
		)
		fmt.Println("Sending: ", response.Type, " to: ", msg.Sender)
		fmt.Println("FIND_NODE_RESPONSE: ", response)
		fmt.Println("Contacts inside resp: ", response.Contacts)
		msgToSend, err := json.Marshal(response)

		return msgToSend, err
	case "FIND_VALUE":
		network.node.RoutingTable.AddContact(*msg.Sender)
		var response KademliaMessage
		data, exists := network.node.storagehandler.getValue(msg.Data.Key)
		if exists {
			response = CreateKademliaMessage(
				"FIND_VALUE_RESPONSE",
				msg.Receiver,
				msg.Sender,
				nil,
				&data,
				nil,
			)

			fmt.Println("Sending: ", response.Type, " to: ", msg.Sender)
			theValueToSend, err := json.Marshal(response)

			return theValueToSend, err
		} else {
			closestNodes := network.node.RoutingTable.FindClosestContacts(msg.Receiver.ID, network.node.k)
			response = CreateKademliaMessage(
				"FIND_VALUE_CONTACTS",
				msg.Receiver,
				msg.Sender,
				msg.Target,
				nil,
				closestNodes,
			)

			fmt.Println("Sending: ", response.Type, " to: ", msg.Sender)
			msgToSend, err := json.Marshal(response)
			return msgToSend, err
		}

	case "STORE":
		var response KademliaMessage
		network.node.RoutingTable.AddContact(*msg.Sender)
		b := network.node.handleStoreMessage(*msg.Data)
		if b {
			response = CreateKademliaMessage(
				"STORE_SUCCESSFUL",
				msg.Receiver,
				msg.Sender,
				nil,
				msg.Data,
				nil,
			)
		} else {
			response = CreateKademliaMessage(
				"STORE_FAILED",
				msg.Receiver,
				msg.Sender,
				nil,
				msg.Data,
				nil,
			)
		}

		fmt.Println("Sending: ", response.Type, " to: ", msg.Sender)
		msgToSend, err := json.Marshal(response)
		return msgToSend, err
	case "FORGET":
		network.node.RoutingTable.AddContact(*msg.Sender)
		b := network.node.handleForgetMessage(*msg.Data)
		var response KademliaMessage
		if b {
			response = CreateKademliaMessage(
				"FORGET_SUCCESSFUL",
				msg.Receiver,
				msg.Sender,
				nil,
				msg.Data,
				nil,
			)
		} else {
			response = CreateKademliaMessage(
				"FORGET_FAILED",
				msg.Receiver,
				msg.Sender,
				nil,
				msg.Data,
				nil,
			)
		}

		fmt.Println("Sending: ", response.Type, " to: ", msg.Sender)
		msgToSend, err := json.Marshal(response)
		return msgToSend, err
	//	HandleStore(msg)
	case "PONG":

	case "FIND_NODE_RESPONSE":

	case "FIND_VALUE_CONTACTS":

	case "FIND_VALUE_RESPONSE":

	case "STORE_SUCCESSFUL":

	case "STORE_FAILED":

	case "FORGET_SUCCESSFUL":

	case "FORGET_FAILED":

	default:

		fmt.Println("Received unknown message type:", msg.Type)
		fmt.Println("Message: ", msg)

	}
	return nil, nil
}

func (network *Network) SendPingMessage(sender *Contact, receiver *Contact) error {

	pingMessage := CreateKademliaMessage(
		"PING",
		sender,
		receiver,
		nil,
		nil,
		[]Contact{},
	)
	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)

	data, err := network.Send(addr, pingMessage)
	if err != nil {
		log.Printf("Ping failed: %v\n", err)
		return err
	}
	var msg KademliaMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("Error decoding message during ping: ", err)
		return err
	}

	fmt.Println(msg.Sender.Address + ":" + strconv.Itoa(msg.Sender.Port) + " AKW YOUR PING")
	return nil
}

func (network *Network) SendFindContactMessage(sender, receiver, target *Contact) ([]Contact, error) {
	fmt.Println("\t\tNode getting the request to find more nodes", receiver.Address+":"+strconv.Itoa(receiver.Port))

	message := CreateKademliaMessage(
		"FIND_NODE",
		sender,
		receiver,
		target,
		nil,
		[]Contact{},
	)
	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)

	data, err := network.Send(addr, message)
	if err != nil {
		log.Printf("FIND_NODE FAILED: %v\n", err)
		return nil, err
	}
	var msg KademliaMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("\t\t\t Error decoding message:", err)
		return nil, err
	}

	return msg.Contacts, nil
}

func (network *Network) SendFindDataMessage(sender, receiver *Contact, hash string) ([]Contact, *StorageData, error) {

	message := CreateKademliaMessage(
		"FIND_VALUE",
		sender,
		receiver,
		nil,
		&StorageData{Key: hash},
		[]Contact{},
	)

	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)

	data, err := network.Send(addr, message)
	if err != nil {
		log.Printf("FIND_NODE FAILED: %v\n", err)
		return nil, nil, err
	}
	var msg KademliaMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("\t\t\t Error decoding message:", err)
		return nil, nil, err
	}

	fmt.Println("\t\t\t msg", msg)
	switch msg.Type {
	case "FIND_VALUE_RESPONSE":
		return []Contact{*msg.Sender}, msg.Data, nil
	case "FIND_VALUE_CONTACTS":
		return msg.Contacts, nil, nil
	default:
		fmt.Println("\t\t\t FIND VALUE GOT A MESSAGE IT DID NOT EXPECT!!!!!!!!!!", msg.Type)
		return nil, nil, nil
	}
}

func (network *Network) SendStoreMessage(sender, receiver *Contact, storageData *StorageData) (string, error) {
	// TODO
	message := CreateKademliaMessage(
		"STORE",
		sender,
		receiver,
		nil,
		storageData,
		nil,
	)

	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)
	fmt.Printf(" %-10s  %-10s  %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10t  \n", "This node: ", message.Sender.Address, " is sending a message to: ", message.Receiver.Address, " Of the type: ", message.Type, " Data to be stored: Key: ", storageData.Key, " value: ", storageData.Value, " ttl: ", storageData.TimeToLive, " original: ", storageData.Original)
	data, err := network.Send(addr, message)
	if err != nil {
		log.Printf("STORE FAILED: %v\n", err)
		return "", err
	}
	var msg KademliaMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("\t\t\t Error decoding message:", err)
		return "", err
	}

	fmt.Println("\t\t\t msg", msg)

	switch msg.Type {
	case "STORE_SUCCESSFUL":
		return msg.Sender.Address + " STORE_SUCCESSFUL", nil
	case "STORE_FAILED":
		return msg.Sender.Address + " STORE_FAILED", err
	default:
		fmt.Println("\t\t\t FIND VALUE GOT A MESSAGE IT DID NOT EXPECT!!!!!!!!!!", msg.Type)
		return "", nil
	}
}

func (network *Network) SendForgetMessage(sender, receiver *Contact, storageData *StorageData) (string, bool, error) {
	// TODO
	message := CreateKademliaMessage(
		"FORGET",
		sender,
		receiver,
		nil,
		storageData,
		nil,
	)

	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)
	data, err := network.Send(addr, message)
	if err != nil {
		log.Printf("FORGET FAILED: %v\n", err)
		return "", false, err
	}
	var msg KademliaMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("\t\t\t Error decoding message during forget command:", err)
		return "", false, err
	}

	switch msg.Type {
	case "FORGET_SUCCESSFUL":
		return msg.Sender.Address + " FORGET_SUCCESSFUL", true, nil
	case "FORGET_FAILED":
		return msg.Sender.Address + " FORGET_SUCCESSFUL", false, err
	default:
		fmt.Println("\t\t\t FORGET VALUE GOT A MESSAGE IT DID NOT EXPECT!!!!!!!!!!", msg.Type)
		return "", false, nil
	}
}

func (network *Network) Send(addr string, msg KademliaMessage) ([]byte, error) {

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return nil, err
	}

	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return nil, err
	}
	defer connection.Close()

	marshaledMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, err
	}
	_, err = connection.Write(marshaledMsg)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return nil, err
	}

	responseChannel := make(chan []byte)
	go func() {
		// Read from the connection
		//data := make([]byte, 1024)
		data := make([]byte, 2024)
		length, _, err := connection.ReadFromUDP(data[:])
		if err != nil {
			return
		}
		responseChannel <- data[:length]

	}()

	select {
	case response := <-responseChannel:
		if _, err := network.Dispatcher(response); err != nil {
			return nil, err
		} else {
			return response, nil
		}
	case <-time.After(3 * time.Second):
		return nil, errors.New("TIME OUT PLS LOOK INTO THIS")
	}
}

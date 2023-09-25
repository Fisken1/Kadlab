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
	Type     string    `json:"Type"`
	Sender   *Contact  `json:"SenderID"`
	Receiver *Contact  `json:"NodeID"`
	Target   *Contact  `json:"TargetID"`
	Key      string    `json:"Key,omitempty"`
	Value    string    `json:"Data,omitempty"`
	Contacts []Contact `json:"contacts,omitempty"`
}

func CreateKademliaMessage(messageType, key, value string, sender, receiver, target *Contact, contacts []Contact) KademliaMessage {
	return KademliaMessage{
		Type:     messageType,
		Sender:   sender,
		Receiver: receiver,
		Target:   target,
		Key:      key,
		Value:    value,
		Contacts: contacts,
	}
}

func (network *Network) Listen(contact Contact) {

	addr := contact.Address + ":" + strconv.Itoa(contact.Port)
	fmt.Println(" Kademlia listener is starting on: " + addr)
	listenAdrs, _ := net.ResolveUDPAddr("udp", addr)
	servr, err := net.ListenUDP("udp", listenAdrs)
	network.server = servr
	if err != nil {
		fmt.Println("BIG ERROR! This is the addr it tried to listen to: "+listenAdrs.String(), err)
		defer servr.Close()
	}
	fmt.Println("Listening on: " + listenAdrs.String() + " " + contact.ID.String() + "\n\n")

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
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return nil, err
	}

	switch msg.Type {
	case "PING":
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return nil, err
		}
		fmt.Println(msg.Sender.Address + ":" + strconv.Itoa(msg.Sender.Port) + " Sent a PING to " + msg.Receiver.Address + ":" + strconv.Itoa(msg.Receiver.Port))
		pong := CreateKademliaMessage(
			"PONG",
			"",
			"",
			msg.Receiver,
			msg.Sender,
			nil,
			[]Contact{},
		)
		msgToSend, err := json.Marshal(pong)
		network.node.RoutingTable.AddContact(*msg.Sender)
		return msgToSend, err

	case "FIND_NODE":
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return nil, err
		}
		fmt.Println(msg.Sender.Address + " <- This node sent a FIND_NODE rpc to us")
		closestNodes := network.node.RoutingTable.FindClosestContacts(msg.Receiver.ID, network.node.k)
		fmt.Println("closest nodes: ", closestNodes)
		find := CreateKademliaMessage(
			"FIND_NODE_RESPONSE",
			"",
			"",
			msg.Receiver,
			msg.Sender,
			msg.Target,
			closestNodes,
		)
		msgToSend, err := json.Marshal(find)
		network.node.RoutingTable.AddContact(*msg.Sender)
		return msgToSend, err
	case "FIND_VALUE":
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return nil, err
		}
		fmt.Println(msg.Sender.Address + " <- This node sent a FIND_VALUE rpc to us")
		fmt.Println("message: ", msg)
		fmt.Println("SENT KEY: ", msg.Key)
		if value, exists := network.node.Hashmap[msg.Key]; exists {
			valueFound := CreateKademliaMessage(
				"FIND_VALUE_RESPONSE",
				msg.Key,
				string(value),
				msg.Receiver,
				msg.Sender,
				nil,
				nil,
			)
			theValueToSend, err := json.Marshal(valueFound)
			fmt.Println("Returning from FIND_VALUE_RESPONSE ")
			return theValueToSend, err
		} else {
			closestNodes := network.node.RoutingTable.FindClosestContacts(msg.Receiver.ID, network.node.k)
			fmt.Println("closest nodes: ", closestNodes)
			find := CreateKademliaMessage(
				"FIND_VALUE_CONTACTS",
				"",
				"",
				msg.Receiver,
				msg.Sender,
				msg.Target,
				closestNodes,
			)
			msgToSend, err := json.Marshal(find)
			fmt.Println("Returning from FIND_VALUE_CONTACTS ")
			//network.node.RoutingTable.AddContact(*msg.Sender)
			return msgToSend, err
		}

	case "STORE":
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return nil, err
		}
		fmt.Println(msg.Sender.Address + " <- This node sent a STORE message to us")
		fmt.Println("GOT A STORE MESSAGE!!!!!!", msg)
		network.node.handleStoreMessage(msg.Key, []byte(msg.Value))
		response := CreateKademliaMessage(
			"STORE_SUCCESSFUL",
			"",
			"",
			msg.Receiver,
			msg.Sender,
			nil,
			nil,
		)
		msgToSend, err := json.Marshal(response)
		return msgToSend, err

	//	HandleStore(msg)
	case "PONG":
		fmt.Println("inside Pong case")
	case "FIND_NODE_RESPONSE":
		fmt.Println("inside FINDE_NODE_RESPONSE case")

	case "STORE_SUCCESSFUL":
		fmt.Println("store response")
	default:
		//network.node.RoutingTable.AddContact(*msg.Sender)
		fmt.Println("Received unknown message typeee:", msg.Type)
		fmt.Println("Sender: ", msg.Sender)
		fmt.Println(msg)
	}
	return nil, nil
}

func (network *Network) SendPingMessage(sender *Contact, receiver *Contact) error {

	pingMessage := CreateKademliaMessage(
		"PING",
		"",
		"",
		sender,
		receiver,
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
		fmt.Println("Error decoding message:", err)
		return err
	}

	fmt.Println(msg.Sender.Address + ":" + strconv.Itoa(msg.Sender.Port) + " AKW YOUR PING")
	return nil
}

func (network *Network) SendFindContactMessage(sender, receiver, target *Contact) ([]Contact, error) {
	fmt.Println("\t\tNode getting the request to find more nodes", receiver.Address+":"+strconv.Itoa(receiver.Port))
	fmt.Println("\t\tTarget we are looking for", target.Address+":"+strconv.Itoa(target.Port))
	message := CreateKademliaMessage(
		"FIND_NODE",
		"",
		"",
		sender,
		receiver,
		target,
		[]Contact{},
	)
	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)
	fmt.Println("\t\tNOW SENDING FIND_NODE")

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
	fmt.Println("\t\t\t len of msg", len(msg.Contacts))
	fmt.Println("\t\t\t msg", msg)
	for _, contact := range msg.Contacts {
		fmt.Println("\t\t\tthis is one contact returned from FIND_NODE rpc", contact)
	}
	return msg.Contacts, nil
}

func (network *Network) SendFindDataMessage(sender, receiver *Contact, hash string) ([]Contact, []byte, error) {
	fmt.Println("\t\tNode getting the request to find more nodes", receiver.Address+":"+strconv.Itoa(receiver.Port))

	message := CreateKademliaMessage(
		"FIND_VALUE",
		hash,
		"",
		sender,
		receiver,
		nil,
		[]Contact{},
	)
	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)
	fmt.Println("\t\tNOW SENDING FIND_VALUE")
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
	fmt.Println("\t\t\t len of msg", len(msg.Contacts))
	fmt.Println("\t\t\t msg", msg)
	switch msg.Type {
	case "FIND_VALUE_RESPONSE":
		return msg.Contacts, []byte(msg.Value), nil
	case "FIND_VALUE_CONTACTS":
		return msg.Contacts, nil, nil
	default:
		fmt.Println("\t\t\t FIND VALUE GOT A MESSAGE IT DID NOT EXPECT!!!!!!!!!!", msg.Type)
		return nil, nil, nil
	}
}

func (network *Network) SendStoreMessage(sender, receiver *Contact, key string, value []byte) (string, error) {
	// TODO
	message := CreateKademliaMessage(
		"STORE",
		key,
		string(value),
		sender,
		receiver,
		nil,
		nil,
	)
	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)
	fmt.Println("\t\tNOW SENDING STORE TO NODE: ", addr)
	data, err := network.Send(addr, message)
	if err != nil {
		log.Printf("FIND_NODE FAILED: %v\n", err)
		return "", err
	}
	var msg KademliaMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("\t\t\t Error decoding message:", err)
		return "", err
	}
	fmt.Println("\t\t\t len of msg", len(msg.Contacts))
	fmt.Println("\t\t\t msg", msg)
	fmt.Println("\t\t\t Attempting to store")
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

func (network *Network) Send(addr string, msg KademliaMessage) ([]byte, error) {
	fmt.Println("we are now sending")
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
		data := make([]byte, 1024)
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

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
	Key      string    `json:"Data,omitempty"`
	Value    string    `json:"Key,omitempty"`
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
		json.Unmarshal(data, &msg)
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
		json.Unmarshal(data, &msg)
		fmt.Println(msg.Sender.Address + " Sent the node printing this a FIND_NODE")
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
	//case "find_value":
	//	HandleFindValue(msg)
	//case "store":
	//	HandleStore(msg)
	default:
		network.node.RoutingTable.AddContact(*msg.Sender)
		fmt.Println("Received unknown message type:", msg.Type)
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

func (network *Network) SendFindDataMessage(sender, receiver, target *Contact, hash string) ([]Contact, error) {
	fmt.Println("\t\tNode getting the request to find more nodes", receiver.Address+":"+strconv.Itoa(receiver.Port))
	fmt.Println("\t\tTarget we are looking for", target.Address+":"+strconv.Itoa(target.Port))
	message := CreateKademliaMessage(
		"FIND_VALUE",
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

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
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

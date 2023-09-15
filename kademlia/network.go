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
	node Kademlia
}

// KademliaMessage represents a Kademlia message.
type KademliaMessage struct {
	Type     string   `json:"Type"`
	Sender   *Contact `json:"SenderID"`
	Receiver *Contact `json:"NodeID"`
	Key      string   `json:"Data,omitempty"`
	Value    string   `json:"Key,omitempty"`
}

func CreateKademliaMessage(messageType, key, value string, sender, receiver *Contact) KademliaMessage {
	return KademliaMessage{
		Type:     messageType,
		Sender:   sender,
		Receiver: receiver,
		Key:      key,
		Value:    value,
	}
}

// This function is taken from "https://go.dev/play/p/BDt3qEQ_2H"
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func (network *Network) Listen(contact Contact) {
	fmt.Println("Kademlia listener is starting...")
	addr := contact.Address + ":" + strconv.Itoa(contact.Port)
	listenAdrs, _ := net.ResolveUDPAddr("udp", addr)

	servr, err := net.ListenUDP("udp", listenAdrs)
	if err != nil {
		fmt.Println("BIG ERROR!!!!!!!!!!!!!!!!!!!!!!!!!!!", err)
		defer servr.Close()
	}
	fmt.Println("Listening on: " + listenAdrs.String() + " " + contact.ID.String() + "\n\n")

	for {
		buf := make([]byte, 65536)
		rlen, _, err := servr.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		// Process the received data by passing it to the Dispatcher function.
		go network.Dispatcher(buf[:rlen])
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
		fmt.Println(msg.Sender.Address + " Sent the node printing this a PING")
		pong := CreateKademliaMessage("PONG", "", "", msg.Receiver, msg.Sender)
		msgToSend, err := json.Marshal(pong)
		return msgToSend, err

	case "FIND_NODE":
		json.Unmarshal(data, &msg)
		fmt.Println(msg.Sender.Address + " Sent the node printing this a FIND_NODE")
		closestNodes := network.node.RoutingTable.FindClosestContacts(msg.Receiver.ID, 7)
		msgToSend, err := json.Marshal(closestNodes)
		return msgToSend, err
	//case "find_value":
	//	HandleFindValue(msg)
	//case "store":
	//	HandleStore(msg)
	default:
		fmt.Println("Received unknown message type:", msg.Type)
	}
	return nil, nil
}

func (network *Network) SendPingMessage(sender *Contact, receiver *Contact) error {
	pingMessage := CreateKademliaMessage("PING", "", "", sender, receiver)
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

	fmt.Println(msg.Receiver.Address + "AKW YOUR PING")
	return nil
}

func (network *Network) SendFindContactMessage(receiver *Contact, target *Contact) ([]Contact, error) {
	pingMessage := CreateKademliaMessage("FIND_NODE", "", "", receiver, target)
	addr := receiver.Address + ":" + strconv.Itoa(receiver.Port)
	data, err := network.Send(addr, pingMessage)
	if err != nil {
		log.Printf("FIND_NODE FAILED: %v\n", err)
		return nil, err
	}
	var msg []Contact
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return nil, err
	}

	return msg, nil
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
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

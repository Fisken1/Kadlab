package kademlia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
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
		go Dispatcher(buf[:rlen])
	}
}

// Dispatcher is responsible for routing incoming messages to their respective handlers.
func Dispatcher(data []byte) {
	var msg KademliaMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}

	switch msg.Type {
	case "PING":
		HandlePing(msg)
	case "PONG":
		HandlePong(msg)
	//case "find_node":
	//	HandleFindNode(msg)
	//case "find_value":
	//	HandleFindValue(msg)
	//case "store":
	//	HandleStore(msg)
	default:
		fmt.Println("Received unknown message type:", msg.Type)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// Create a UDP connection to the contact.
	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer connection.Close()

	pingMessage := CreateKademliaMessage("PING", contact.ID.String(), contact.Address, "", "", "")

	msg, err := json.Marshal(pingMessage)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	_, err = connection.Write(msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

// BELOW WE HANDLE RPCs

// HandlePing handles a "ping" message.
func HandlePing(msg KademliaMessage) {
	fmt.Println(msg.Receiver.Address+":"+strconv.Itoa(msg.Receiver.Port)+" Received "+msg.Type+" from", msg.Sender.Address+":"+strconv.Itoa(msg.Receiver.Port))
	msg.Sender.SendPongMessage(msg.Receiver)
}

// HandlePong handles a "pong" message.
func HandlePong(msg KademliaMessage) {
	fmt.Println(msg.Receiver.Address+":"+strconv.Itoa(msg.Receiver.Port)+" Received "+msg.Type+" from", msg.Sender.Address+":"+strconv.Itoa(msg.Receiver.Port))
}

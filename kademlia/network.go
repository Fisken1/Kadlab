package kademlia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Network struct {
	Contact  *Contact
	Kademlia *Kademlia
}

// KademliaMessage represents a Kademlia message.
type KademliaMessage struct {
	Type          string `json:"Type"`
	SenderID      string `json:"SenderID"`
	SenderAddress string `json:"SenderAddress"`
	NodeID        string `json:"NodeID"`
	Data          string `json:"Data,omitempty"`
	Key           string `json:"Key,omitempty"`
}

func CreateKademliaMessage(messageType, senderID, senderAddress, targetNodeID, data, key string) KademliaMessage {
	return KademliaMessage{
		Type:          messageType,
		SenderID:      senderID,
		SenderAddress: senderAddress,
		NodeID:        targetNodeID,
		Data:          data,
		Key:           key,
	}
}

func InitNode(me *Contact) *Network {
	network := &Network{
		Contact:  me,
		Kademlia: InitKademila(*me),
	}
	fmt.Print("INITNODE", me.distance, me.ID, me.Address)
	return network
}

func InitBootstrap(port int) *Network {
	ip, err := externalIP()
	if err != nil {
		fmt.Print(err)
	}
	me := NewContact(NewRandomKademliaID(), ip, port)
	node := InitNode(&me)
	return node
}

func InitJoin(bootstrapIP string, port int) (*Network, error) {
	// Step 1: Get the local IP address from eth0.
	ip, err := externalIP()
	if err != nil {
		return nil, err
	}
	n := NewRandomKademliaID()

	// Step 3: Create a new contact for the bootstrap node.
	bootstrapContact := NewContact(nil, ip, port)

	// Step 4: Initialize the Kademlia instance for the new node.
	node := NewContact(n, ip, port)

	newNodeNetwork := InitNode(&node)

	newNodeNetwork.Kademlia.RoutingTable.AddContact(bootstrapContact)

	return newNodeNetwork, nil
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

func Listen(ip string, port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	listener, err := net.ListenUDP("udp", &addr)

	if err != nil {
		panic(err)
	}
	defer listener.Close()

	buf := make([]byte, 1024)

	for {
		rlen, _, err := listener.ReadFromUDP(buf)
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
	fmt.Println("Received ping from", msg.SenderID)
	// Handle the ping message logic here.
}

// HandlePong handles a "pong" message.
func HandlePong(msg KademliaMessage) {
	fmt.Println("Received pong from", msg.SenderID)
	// Handle the pong message logic here.
}

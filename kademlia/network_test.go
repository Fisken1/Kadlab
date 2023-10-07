package kademlia

/*
import (
	"log"
	"net"
	"testing"
	"time"
)

func TestNetworkSendPingMessage(t *testing.T) {
	ip := "192.168.1.26"
	ipBootstrap := GetBootstrapIP(ip) // Same IP as the current node for bootstrap
	bootstrap, err := InitJoin(ipBootstrap, 5000)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Use a channel for synchronization.
	done := make(chan bool)

	// Start a goroutine to close the network and signal when done.
	go func() {
		// Sleep for a while to allow the network to start.
		time.Sleep(1 * time.Second)

		done <- true
	}()
	bootstrap.fixNetwork()
	// Wait for the goroutine to finish before proceeding with other checks.
	<-done

	ip = "192.168.1.26"
	port := 57707
	node, err := InitJoin(ip, port)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Use a channel for synchronization.
	done = make(chan bool)

	// Start a goroutine to close the network and signal when done.
	go func() {
		// Sleep for a while to allow the network to start.
		time.Sleep(1 * time.Second)

		done <- true
	}()

	// Wait for the goroutine to finish before proceeding with other checks.
	<-done
	node.fixNetwork()
	// Call the SendPingMessage method
	lookatthiserr := bootstrap.net.SendPingMessage(&bootstrap.RoutingTable.me, &node.RoutingTable.me)

	// Check for errors
	if lookatthiserr != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	// Add more specific checks based on your application logic
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}


*/

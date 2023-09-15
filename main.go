package main

import "Kadlab/kademlia"

func main() {
	node, _ := kademlia.InitJoin(8080)
	kademlia.Cli(node, 9090)
	select {}
}

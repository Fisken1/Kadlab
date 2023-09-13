package main

import "Kadlab/kademlia"

func main() {
	node, _ := kademlia.InitJoin(8080)

	go kademlia.Cli(node, make(chan int))
	select {}

}

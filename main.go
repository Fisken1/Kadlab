package main

import "Kadlab/kademlia"

func main() {
	kademlia.InitJoin(8080)
	select {}
}

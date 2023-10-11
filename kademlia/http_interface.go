package kademlia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type HTTPInterface struct {
	kademlia     *Kademlia
	routingTable *RoutingTable
}

var httpInterF *HTTPInterface

func InitHTTPInterface(kad *Kademlia, rT *RoutingTable) {
	var httpIF HTTPInterface = HTTPInterface{
		kademlia:     kad,
		routingTable: rT,
	}
	httpInterF = &httpIF

	http.HandleFunc("/objects/", handleReq)
	http.ListenAndServe(":5000", nil)

}

func handleReq(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		//println("req.Uri: ", req.RequestURI)
		s := strings.Split(req.RequestURI, "/objects/")
		//println("hash: ", s[1])
		b := []byte(s[1])
		hash, location := httpInterF.kademlia.Store(b)
		fmt.Println("hash from Store: ", hash, " at location: ", location)

		if hash != "0" {
			//resp = location + "/objects/" + hash

			fmt.Println("Stored")
		} else {
			fmt.Println("store failed")
			//resp = "Error, failed to store..."
		}
		jsonResp, err := json.Marshal(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		//w.Header().Set("Location", resp)

		w.Write(jsonResp)
	}
	if req.Method == "GET" {
		fmt.Println("req: ", req.RequestURI)

		//Trim req.RequestURI so only the sent hash remains
		hash := strings.Split(req.RequestURI, "/objects/")[1]

		fmt.Println("hash: ", hash)
		response := ""
		_, _, data, err := httpInterF.kademlia.LookupData(hash)
		if err != nil {
			fmt.Println(err)

		} else {
			if &data != nil {
				response = string(data.Value)

			}
		}
		jsonResp, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(jsonResp)

	}

}

package kademlia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	http.HandleFunc("/objects", handlePost)
	http.HandleFunc("/objects/", handleGet)
	http.ListenAndServe(":1550", nil)
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	println(req.RemoteAddr)
	fmt.Println("Post")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {

	}
	fmt.Println("body: ", body)
	hash, location := httpInterF.kademlia.Store(body)
	resp := ""
	if hash != "0" {
		resp = location + "/objects/{" + hash + "}"

		fmt.Println("Stored")
	} else {
		resp = "Error, failed to store..."
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", resp)

	w.Write(nil)
}
func handleGet(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Get")

	fmt.Println("req: ", req.RequestURI)

	//Trim req.RequestURI so only the sent hash remains
	resString := strings.ReplaceAll(req.RequestURI, "%7B", " ")
	resString = strings.ReplaceAll(resString, "%7D", " ")
	s := strings.Split(resString, " ")
	if len(s) != 3 {
		fmt.Println("error in get s: ", s)
		return
	}
	hash := s[1]

	fmt.Println("hash: ", hash)
	response := ""
	_, _, data, err := httpInterF.kademlia.LookupData(hash)
	if err != nil {
		fmt.Println(err)

	} else {
		if data.Value != nil {
			response = string(data.Value)

		}
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(jsonResp)

	w.Write(nil)

}

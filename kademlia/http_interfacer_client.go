package kademlia

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func httpClient() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("httpClient> ")
	for {
		scanner.Scan()
		text := scanner.Text()

		if len(text) > 0 {
			input := strings.Fields(text)
			answer := httpClientInputHandler(input)
			fmt.Print(answer + "httpClient> ")
			//TODO
		} else {
			fmt.Print("httpClient> ")
		}
	}
}
func httpClientInputHandler(input []string) string {

	switch input[0] {
	case "POST":
		fmt.Println("Post")
		address := "http://" + input[1] + "/objects"
		fmt.Println("address ", address)
		inputString := input[2]
		fmt.Println("inputStrings: ", inputString)
		body, _ := json.Marshal(inputString)

		//resp, err := http.Post(address, "application/json", bytes.NewBuffer(body))
		_, err := http.Post(address, "application/json", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		//defer resp.Body.Close()
		return ""
	case "GET":
		address := "http://" + input[1] + "/objects/{testing12345}}}"
		fmt.Println("address ", address)
		//body, _ := json.Marshal(inputString)

		//resp, err := http.Post(address, "application/json", bytes.NewBuffer(body))
		resp, err := http.Get(address)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		return ""

	default:

	}

	return ""
}

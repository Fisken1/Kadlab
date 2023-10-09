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
			//fmt.Println(input)
		} else {
			fmt.Print("httpClient> ")
		}
	}
}
func httpClientInputHandler(input []string) string {

	switch input[0] {
	case "POST":
		fmt.Println("Post")
		address := "http://localhost:" + input[1] + "/objects"
		fmt.Println("address ", address)
		inputString := input[2]
		fmt.Println("inputStrings: ", inputString)
		body, _ := json.Marshal(inputString)

		_, err := http.Post(address, "application/json", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		return ""
	case "GET":
		address := "http://localhost:" + input[1] + "/objects/{" + input[2] + "}"
		fmt.Println("address ", address)
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

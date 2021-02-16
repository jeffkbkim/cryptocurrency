package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jeffkbkim/cryptocurrency/gossip"
)

func GetBalance(args ...string) {
	addr := args[0]
	user := args[1]
	url := fmt.Sprintf("http://localhost:%s/balance/%s", addr, user)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}

func CreateUser(args ...string) {
	addr := args[0]
	user := args[1]
	url := fmt.Sprintf("http://localhost:%s/users/%s", addr, user)
	resp, err := http.Post(url, "json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}

func CreateTransfer(args ...string) {
	addr := args[0]
	to := args[1]
	amount := args[2]
	url := fmt.Sprintf("http://localhost:%s/transfers/%s/%s", addr, to, amount)
	resp, err := http.Post(url, "json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}

func Gossip(port string, state *gossip.State) (*gossip.State, int) {
	url := fmt.Sprintf("http://localhost:%s/gossip", port)
	requestBody, err := json.Marshal(*state)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
		return nil, http.StatusNotFound
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var theirState gossip.State

	err = decoder.Decode(&theirState)
	if err != nil {
		log.Fatal(err)
	}

	return &theirState, resp.StatusCode
}

func PubKey(port string) string {
	url := fmt.Sprintf("http://localhost:%s/pub_key", port)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes)
}

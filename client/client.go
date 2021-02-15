package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	from := args[1]
	to := args[2]
	amount := args[3]
	url := fmt.Sprintf("http://localhost:%s/transfers/%s/%s/%s", addr, from, to, amount)
	resp, err := http.Post(url, "json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}

func Gossip(peer string, state map[string]gossip.Pair) (map[string]gossip.Pair, int) {
	url := fmt.Sprintf("http://localhost:%s/gossip", peer)
	requestBody, err := json.Marshal(state)
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
	var theirState map[string]gossip.Pair

	err = decoder.Decode(&theirState)
	if err != nil {
		log.Fatal(err)
	}

	return theirState, resp.StatusCode
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	functions = map[string]func(){
		"getBalance":     getBalance,
		"createUser":     createUser,
		"createTransfer": createTransfer,
	}

	addr   = "8080"
	method = flag.String("func", "", "name of function to call")
	user   = flag.String("user", "", "")
	from   = flag.String("from", "", "")
	to     = flag.String("to", "", "")
	amount = flag.String("amount", "", "")
)

func main() {
	flag.Parse()

	function, ok := functions[*method]
	if !ok {
		log.Fatal("function does not exist")
	}

	function()
}

func getBalance() {
	url := fmt.Sprintf("http://localhost:%s/balance/%s", addr, *user)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}

func createUser() {
	url := fmt.Sprintf("http://localhost:%s/users/%s", addr, *user)
	resp, err := http.Post(url, "json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}

func createTransfer() {
	url := fmt.Sprintf("http://localhost:%s/transfers/%s/%s/%s", addr, *from, *to, *amount)
	resp, err := http.Post(url, "json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}

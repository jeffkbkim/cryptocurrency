package main

import (
	"flag"
	"fmt"
	"log"

	client "github.com/jeffkbkim/cryptocurrency/client"
)

var (
	functions = map[string]func(args ...string){
		"getBalance":     client.GetBalance,
		"createUser":     client.CreateUser,
		"createTransfer": client.CreateTransfer,
	}
)

func main() {
	flag.Parse()
	method := flag.Args()[0]
	function, ok := functions[method]
	if !ok {
		fmt.Println(method)
		log.Fatal("function does not exist")
	}

	function(flag.Args()[1:]...)
}

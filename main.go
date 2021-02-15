package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	BALANCES = map[string]int{"binpy": 1_000_000}
)

func main() {
	addr := flag.String("port", "8080", "port to open")
	peer := flag.String("peer", "", "peer port to connect to")

	log.Fatal(run(*addr, *peer))
}

func run(addr string, peer string) error {
	mux := makeMuxRouter()
	log.Println("Listening on", addr)
	s := &http.Server{
		Addr:           ":" + addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/balance/{id}", handleGetBalance).Methods("GET")
	muxRouter.HandleFunc("/users/{id}", handleCreateUser).Methods("POST")
	muxRouter.HandleFunc("/transfers/{from}/{to}/{amount}", handleCreateTransfer).Methods("POST")
	return muxRouter
}

func handleGetBalance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, ok := BALANCES[id]; !ok {
		fmt.Println(id, "does not exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("%s has %d\n", id, BALANCES[id])
	io.WriteString(w, strconv.Itoa(BALANCES[id]))
	printBalances()
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if BALANCES[id] != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	BALANCES[id] = 0
	fmt.Println("OK")
	io.WriteString(w, "OK")
	printBalances()
}

func handleCreateTransfer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	from := params["from"]
	to := params["to"]
	amountStr := params["amount"]
	amount, err := strconv.Atoi(amountStr)
	fmt.Println("transfering amount", amount)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	if _, ok := BALANCES[to]; !ok {
		fmt.Println(to, "does not exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if BALANCES[from] < amount {
		fmt.Println(from, "has insufficient funds", BALANCES[from])
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	BALANCES[from] -= amount
	BALANCES[to] += amount
	fmt.Println("OK")
	io.WriteString(w, "OK")
	printBalances()
}

func printBalances() {
	b, err := json.MarshalIndent(BALANCES, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

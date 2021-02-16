package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/jeffkbkim/cryptocurrency/client"
	"github.com/jeffkbkim/cryptocurrency/gossip"
	"github.com/jeffkbkim/cryptocurrency/pki"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"

	bc "github.com/jeffkbkim/cryptocurrency/blockchain"
)

var (
	mu    = &sync.Mutex{}
	STATE = &gossip.State{
		Peers: map[string]bool{},
	}
)

func main() {
	flag.Parse()
	port := flag.Args()[0]
	var peer string

	STATE.Blockchain = &bc.Blockchain{
		Blocks: []*bc.Block{},
	}
	privateKey := pki.GenerateKeyPair()
	pubdata, privdata := pki.GetPubPriv(privateKey)
	STATE.PublicKey = string(pubdata)
	// fmt.Println("my pub key", STATE.PublicKey)
	STATE.PrivateKey = string(privdata)
	STATE.Me = port
	STATE.Peers[port] = true

	if len(flag.Args()) > 1 {
		peer = flag.Args()[1]
		STATE.Peers[peer] = true
	} else {
		// you are the progenitor
		STATE.Blockchain.CreateGenesisBlock(STATE.PublicKey, STATE.PrivateKey)
	}

	go func() {
		for {
			for peer := range STATE.Peers {
				if peer == port {
					continue
				}
				fmt.Println("Gossiping with...", peer)

				mu.Lock()
				copyState := gossip.State{}
				err := copier.Copy(&copyState, STATE)
				mu.Unlock()
				if err != nil {
					panic(err)
				}
				theirState, statusCode := client.Gossip(peer, &copyState)

				if statusCode != http.StatusOK {
					mu.Lock()
					delete(STATE.Peers, peer)
					color.Red("%s has disconnected from the network :(", peer)
					mu.Unlock()
				}
				STATE.UpdateState(mu, theirState)
			}

			STATE.RenderState()
			time.Sleep(3 * time.Second)
		}
	}()

	log.Fatal(run(port, peer))
}

func run(addr string, peer string) error {
	mux := makeMuxRouter()
	log.Println("Listening on", addr)
	s := &http.Server{
		Addr:           ":" + addr,
		Handler:        mux,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/gossip", handleGossip).Methods("POST")
	muxRouter.HandleFunc("/transfers/{to}/{amount}", handleCreateTransfer).Methods("POST")
	muxRouter.HandleFunc("/pub_key", handlePubKey).Methods("GET")
	return muxRouter
}

func handleGossip(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var theirState gossip.State

	err := decoder.Decode(&theirState)
	if err != nil {
		log.Fatal(err)
	}

	STATE.UpdateState(mu, &theirState)

	mu.Lock()
	respBody, err := json.Marshal(STATE)
	mu.Unlock()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCreateTransfer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	toPort := params["to"]
	to := client.PubKey(toPort)
	// fmt.Println("to pubKey:", to)
	amountStr := params["amount"]
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	fmt.Println("transferring amount", amount)
	txn := &bc.Transaction{
		From:   STATE.PublicKey,
		To:     to,
		Amount: amount,
	}
	txn.Signature = pki.Sign(txn.ToMessage(), STATE.PrivateKey)
	STATE.Blockchain.AddToChain(txn)
}

func handlePubKey(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(STATE.PublicKey))
	if err != nil {
		panic(err)
	}
}

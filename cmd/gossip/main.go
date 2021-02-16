package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"

	"github.com/jeffkbkim/cryptocurrency/gossip"
)

var (
	mu    = sync.Mutex{}
	STATE = map[string]gossip.Pair{}
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	movies, err := readLines(path + "/movies.txt")
	if err != nil {
		log.Fatal(err)
	}
	movie := getRandomElement(movies)
	version := 0
	color.Green("my favorite movie, now and forever, is %s %d\n", movie, version)

	currState := map[string]gossip.Pair{
		port: gossip.Pair{
			Movie:   movie,
			Version: version,
		},
	}
	if peer != "" {
		currState[peer] = gossip.Pair{}
	}

	updateState(currState)

	go func() {
		for {
			time.Sleep(8 * time.Second)
			color.Yellow("You know what, screw %s, it's so cliche.", movie)
			movie = getRandomElement(movies)
			version++
			color.Green("My new favorite movie is %s %d", movie, version)

			currState = map[string]gossip.Pair{
				port: gossip.Pair{
					Movie:   movie,
					Version: version,
				},
			}

			updateState(currState)
		}
	}()

	// go func() {
	// 	for {
	// 		copyState := map[string]gossip.Pair{}
	// 		mu.Lock()
	// 		for k, v := range STATE {
	// 			copyState[k] = v
	// 		}
	// 		mu.Unlock()
	// 		for peer := range copyState {
	// 			if peer == port {
	// 				continue
	// 			}
	// 			fmt.Println("Gossiping with...", peer)

	// 			theirState, statusCode := client.Gossip(peer, copyState)
	// 			updateState(theirState)
	// 			if statusCode != http.StatusOK {
	// 				mu.Lock()
	// 				delete(STATE, peer)
	// 				color.Red("%s has disconnected from the network :(", peer)
	// 				mu.Unlock()
	// 			}
	// 		}

	// 		renderState()
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()

	log.Fatal(run(port))
}

func run(addr string) error {
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

func updateState(state map[string]gossip.Pair) {
	mu.Lock()
	defer mu.Unlock()
	for port, newPair := range state {
		if port == "" {
			continue
		}
		currPair, ok := STATE[port]
		if ok && newPair.Version <= currPair.Version {
			continue
		}
		// node has disconnected from network
		if !ok && newPair.Version > 0 {
			continue
		}
		STATE[port] = newPair
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getRandomElement(slice []string) string {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	return slice[r.Intn(len(slice))]
}

func renderState() {
	mu.Lock()
	defer mu.Unlock()
	for port, pair := range STATE {
		fmt.Println("-----------------------------------------------")
		fmt.Println(port, "currently likes", pair.Movie, pair.Version)
		fmt.Println("-----------------------------------------------")
	}
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/gossip", handleGossip).Methods("POST")
	return muxRouter
}

func handleGossip(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var theirState map[string]gossip.Pair

	err := decoder.Decode(&theirState)
	if err != nil {
		log.Fatal(err)
	}

	updateState(theirState)

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

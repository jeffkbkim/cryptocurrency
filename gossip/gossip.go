package gossip

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
	bc "github.com/jeffkbkim/cryptocurrency/blockchain"
	"github.com/jeffkbkim/cryptocurrency/helper"
)

type State struct {
	Blockchain *bc.Blockchain
	Me         string
	Peers      map[string]bool
	PublicKey  string
	PrivateKey string
}

func (s *State) UpdateState(mu *sync.Mutex, state *State) {
	mu.Lock()
	defer mu.Unlock()
	if state == nil {
		return
	}

	s.UpdatePeers(state.Peers)
	s.UpdateBlockchain(state.Blockchain)
}

func (s *State) UpdatePeers(peers map[string]bool) {
	for peer := range peers {
		s.Peers[peer] = true
	}
}

func (s *State) UpdateBlockchain(blockchain *bc.Blockchain) {
	if len(s.Blockchain.Blocks) >= len(blockchain.Blocks) {
		return
	}
	if !blockchain.IsValid() {
		return
	}
	s.Blockchain = blockchain
}

func (s *State) RenderState() {
	if s.Blockchain == nil {
		return
	}
	color.Green("              -----------------------------------------------")
	color.Green("-----------------------------Rendering State-----------------------------------------")
	fmt.Println("My blockchain")
	s.Blockchain.Print()
	fmt.Println("blockchain length:", len(s.Blockchain.Blocks))
	fmt.Println("my port", s.Me)
	color.Cyan("my human readable name: %s", helper.HumanReadableName(s.PublicKey))
	color.Yellow("my peers: %v\n", s.Peers)
	s.PrintReadableBalances()
	color.Green("-------------------------------------------------------------------------------------")
	color.Green("              -----------------------------------------------")
}

func (s *State) PrintReadableBalances() {
	for pubKey, balance := range s.Blockchain.ComputeBalances() {
		color.Magenta("%s currently has %d\n", helper.HumanReadableName(pubKey), balance)
	}
}

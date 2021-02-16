package blockchain

import (
	"fmt"

	"github.com/jeffkbkim/cryptocurrency/pki"

	"github.com/fatih/color"

	"github.com/jeffkbkim/cryptocurrency/pow"
)

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	PrevBlockHash string
	Transaction   *Transaction
	Hash          string
	Nonce         string
	Count         int
}

func (b *Blockchain) AddToChain(txn *Transaction) {
	newBlock := &Block{
		Transaction: txn,
	}
	if len(b.Blocks) > 0 {
		newBlock.PrevBlockHash = b.Blocks[len(b.Blocks)-1].Hash
	}
	newBlock.MineBlock()
	b.Blocks = append(b.Blocks, newBlock)
	fmt.Println(newBlock.PrevBlockHash)
	newBlock.prettyPrint()
}

func (b *Blockchain) IsValid() bool {
	for i, block := range b.Blocks {
		if !block.isValid() {
			return false
		}
		if i == len(b.Blocks)-1 {
			continue
		}
		nextBlock := b.Blocks[i+1]
		if block.Hash != nextBlock.PrevBlockHash {
			return false
		}
	}
	return b.isAllSpendsValid()
}

func (b *Blockchain) isAllSpendsValid() bool {
	balances := b.ComputeBalances()
	for _, balance := range balances {
		if balance < 0 {
			return false
		}
	}
	return true
}

func (b *Blockchain) ComputeBalances() map[string]int {
	balances := map[string]int{}
	genesisTxn := b.Blocks[0].Transaction
	balances[genesisTxn.To] = genesisTxn.Amount

	for i, block := range b.Blocks {
		if i == 0 {
			continue
		}
		txn := block.Transaction
		balances[txn.From] -= txn.Amount
		balances[txn.To] += txn.Amount
	}
	return balances
}

func (b *Blockchain) CreateGenesisBlock(pubKey string, privKey string) {
	genesisTxn := &Transaction{
		From:   "",
		To:     pubKey,
		Amount: 500_000,
	}
	genesisTxn.Signature = pki.Sign(genesisTxn.ToMessage(), privKey)
	b.AddToChain(genesisTxn)
}

func (b *Blockchain) Print() {
	if b.Blocks == nil {
		return
	}
	for _, block := range b.Blocks {
		color.Yellow("PrevBlockHash: %s, Own hash: %s, Nonce: %s, Txn: %s\n", block.PrevBlockHash, block.Hash, block.Nonce, block.Transaction.ToMessage())
	}
}

func (b *Block) MineBlock() {
	result := pow.FindNonce(b.blockContents())
	b.Nonce = result.Nonce
	b.Hash = result.Hash
	b.Count = result.Count
}

func (b *Block) isValid() bool {
	_, ok := pow.IsValidNonce(b.blockContents(), b.Nonce)
	return ok && b.Transaction.IsValidSignature()
}

func (b *Block) blockContents() string {
	return b.PrevBlockHash + b.Transaction.ToMessage()
}

func (b *Block) prettyPrint() {
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------------------------------")
	color.Yellow("      Previous hash: %s\n", b.PrevBlockHash)
	color.Green("            Message: %s\n", b.Transaction.ToMessage())
	color.Red("              Nonce: %s\n", b.Nonce)
	color.Yellow("           Own hash: %s\n", b.Hash)
	color.Cyan("              Count: %d\n", b.Count)
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Println("                                         |")
	fmt.Println("                                         |")
	fmt.Println("                                         |")
	fmt.Println("                                         |")
	fmt.Println("                                         |")
	fmt.Println("                                         |")
	fmt.Println("                                         |")
	fmt.Println("                                         |")
	fmt.Printf("                                         %c\n", '\u2193')
}

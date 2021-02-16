package blockchain

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/jeffkbkim/cryptocurrency/pow"
)

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	PrevBlockHash string
	Message       string
	Hash          string
	Nonce         string
	Count         int
}

func (b *Blockchain) AddToChain(message string) {
	newBlock := &Block{
		Message: message,
	}
	if len(b.Blocks) > 0 {
		newBlock.PrevBlockHash = b.Blocks[len(b.Blocks)-1].Hash
	}
	newBlock.MineBlock()
	b.Blocks = append(b.Blocks, newBlock)

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
	return true
}

func (b *Blockchain) CreateGenesisBlock(message string) {
	b.AddToChain(message)
}

func (b *Block) MineBlock() {
	result := pow.FindNonce(b.blockContents())
	b.Nonce = result.Nonce
	b.Hash = result.Hash
	b.Count = result.Count
}

func (b *Block) isValid() bool {
	_, ok := pow.IsValidNonce(b.blockContents(), b.Nonce)
	return ok
}

func (b *Block) blockContents() string {
	return b.PrevBlockHash + b.Message
}

func (b *Block) prettyPrint() {
	// fmt.Println(text.Pad("Previous hash: "+b.PrevBlockHash, 20, '-'))
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------------------------------")
	color.Yellow("      Previous hash: %s\n", b.PrevBlockHash)
	color.Green("            Message: %s\n", b.Message)
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

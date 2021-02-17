package main

import (
	"fmt"

	"github.com/jeffkbkim/cryptocurrency/pki"
	"github.com/jeffkbkim/cryptocurrency/pow"
)

func main() {
	privateKey := pki.GenerateKeyPair()
	pubdata, privdata := pki.GetPubPriv(privateKey)
	signature := pki.Sign("Hello World", string(privdata))
	isValid := pki.Verify("Hello World", signature, string(pubdata))
	fmt.Println(isValid)
	result := pow.FindNonce("Hello World")
	fmt.Println()
	fmt.Println("Nonce:", result.Nonce, "Hash:", result.Hash, "Count:", result.Count)
}

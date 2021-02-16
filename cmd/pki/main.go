package main

import (
	"fmt"

	"github.com/jeffkbkim/cryptocurrency/pki"
	"github.com/jeffkbkim/cryptocurrency/pow"
)

func main() {
	privateKey := pki.GenerateKeyPair()

	pubdata, privdata := pki.GetPubPriv(privateKey)

	// fmt.Println(string(privdata))

	signature := pki.Sign("Hello World", string(privdata))

	// fmt.Println(base64.StdEncoding.EncodeToString(signature))

	isValid := pki.Verify("Hello World", signature, string(pubdata))

	fmt.Println(isValid)

	result := pow.FindNonce("Hello World")
	fmt.Println()
	fmt.Println("Nonce:", result.Nonce, "Hash:", result.Hash, "Count:", result.Count)
}

package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/jeffkbkim/cryptocurrency/pki"
	"github.com/jeffkbkim/cryptocurrency/pow"
)

func main() {
	privateKey := pki.GenerateKeyPair()

	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	fmt.Println(string(pemData))

	signature := pki.Sign("Hello World", privateKey)

	fmt.Println(base64.StdEncoding.EncodeToString(signature))

	isValid := pki.Verify("Hello World", signature, &privateKey.PublicKey)

	fmt.Println(isValid)

	result := pow.FindNonce("Hello World")
	fmt.Println()
	fmt.Println("Nonce:", result.Nonce, "Hash:", result.Hash, "Count:", result.Count)
}

package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

func GenerateKeyPair() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

func Sign(plaintext string, privdata string) []byte {
	block, _ := pem.Decode([]byte(privdata))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash([]byte(plaintext)))
	if err != nil {
		panic(err)
	}

	return signature
}

func Verify(ciphertext string, signature []byte, pubdata string) bool {
	block, _ := pem.Decode([]byte(pubdata))
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash([]byte(ciphertext)), signature)
	if err != nil {
		fmt.Println("verify failed.")
		log.Fatal(err)
	}
	return err == nil
}

func hash(data []byte) []byte {
	s256 := sha256.Sum256(data)
	return s256[:]
}

func GetPubPriv(privateKey *rsa.PrivateKey) ([]byte, []byte) {
	privdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	pubdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		},
	)

	return pubdata, privdata
}

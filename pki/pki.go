package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
)

func GenerateKeyPair() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

func Sign(plaintext string, privateKey *rsa.PrivateKey) []byte {
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash([]byte(plaintext)))
	if err != nil {
		panic(err)
	}

	return signature
}

func Verify(ciphertext string, signature []byte, publicKey *rsa.PublicKey) bool {
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash([]byte(ciphertext)), signature)
	if err != nil {
		log.Fatal(err)
	}
	return err == nil
}

func hash(data []byte) []byte {
	s256 := sha256.Sum256(data)
	return s256[:]
}

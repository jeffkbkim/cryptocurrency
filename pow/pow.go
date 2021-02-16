package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	DIFFICULTY = 5
)

type FindNonceResult struct {
	Nonce string
	Hash  string
	Count int
}

func Hash(message string) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func FindNonce(content string) FindNonceResult {
	nonce := "my first blockchainz"
	count := 0

	var ok bool
	var hash string
	for {
		if count%100_000 == 0 {
			fmt.Printf(".")
		}
		hash, ok = IsValidNonce(content, nonce)
		if ok {
			break
		}
		nonce = NextLexicographicalWord(nonce)
		count++
	}
	fmt.Println()
	return FindNonceResult{
		Nonce: nonce,
		Hash:  hash,
		Count: count,
	}
}

func IsValidNonce(content string, nonce string) (string, bool) {
	threshold := make([]byte, DIFFICULTY)
	for i := 0; i < DIFFICULTY; i++ {
		threshold[i] = '0'
	}
	hash := Hash(content + nonce)
	return hash, strings.HasPrefix(hash, string(threshold))
}

func NextLexicographicalWord(word string) string {
	byteSlice := []byte(word)

	// increment from far back. if digit exceeds 'z', increment next digit.
	i := len(byteSlice) - 1
	for ; i >= 0; i-- {
		byteSlice[i] = byteSlice[i] + 1
		if byteSlice[i] < 'z' {
			break
		}
		byteSlice[i] = 'a'
	}
	if i == -1 {
		byteSlice = append([]byte{'a'}, byteSlice...)
	}

	return string(byteSlice)
}

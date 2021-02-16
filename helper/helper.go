package helper

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"os"
	"strconv"
)

var (
	NAMES, _ = readLines("/Users/jeff.kim/cryptocurrency/helper/human-readable-names.txt")
)

func HumanReadableName(pubKey string) string {
	s256 := sha256.Sum256([]byte(pubKey))
	hash := s256[:]
	hexStr := hex.EncodeToString(hash)

	hashInt := new(big.Int)
	hashInt.SetString(hexStr, 16)
	len := big.NewInt(int64(len(NAMES)))
	mod := new(big.Int).Mod(hashInt, len)

	result, err := strconv.Atoi(mod.String())
	if err != nil {
		panic(err)
	}
	return NAMES[result]
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

func Hash(message string) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/mr-tron/base58"
)

func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex[2:], 16)

	return n
}
func main() {
	log.Println(hexToBigInt("0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"))
	encoded := "HHbNm24TvgRz4ayUg9nxwBWXaw6E74VYJTS1F84D4mN"
	decoded, _ := base58.Decode(encoded)
	decodedStr := hex.EncodeToString(decoded)
	hexStr := fmt.Sprintf("0x%s", decodedStr)
	log.Println(decodedStr)
	// decodedStr = fmt.Sprint("0x%s", decodedStr)
	log.Println("______", hexToBigInt(hexStr))
}

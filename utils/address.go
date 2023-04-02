package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
)

func Base58ToBigIntString(input string) (ret string, err error) {

	decoded, _ := base58.Decode(input)
	decodedStr := hex.EncodeToString(decoded)
	hexStr := fmt.Sprintf("0x%s", decodedStr)
	ret = hexToBigInt(hexStr).String()
	return
}
func Base58ToHexString(input string) (ret string, err error) {
	decoded, _ := base58.Decode(input)
	decodedStr := hex.EncodeToString(decoded)
	hexStr := fmt.Sprintf("0x%s", decodedStr)
	ret = hexStr
	return
}
func HexNumberToBigIntString(input string) (ret string, err error) {
	if !strings.HasPrefix(input, "0x") {
		err = fmt.Errorf("错误的Input 格式%s", input)
		return
	}
	ret = hexToBigInt(input).String()
	return
}

func GetUniqAddress(address string, chainType string) (ret string, err error) {
	if strings.HasPrefix(address, "0x") {
		ret, err = HexNumberToBigIntString(address)
		if err != nil {
			return
		}
		return
	}
	if chainType == "near" {
		ret, err = Base58ToBigIntString(address)
		if err != nil {
			return
		}
		return
	}
	err = fmt.Errorf("没有适配的addres converter")
	return
}

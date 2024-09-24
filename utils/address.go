package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
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
		err = fmt.Errorf("invalid input format %s", input)
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
	if chainType == "solana" {
		ret, err = Base58ToBigIntString(address)
		if err != nil {
			return
		}
		return
	}
	err = fmt.Errorf("no matching address converter")
	return
}
func GetHexAddress(address string, evmType string) (ret string, err error) {
	ret = ""
	if evmType == "near" {
		tokenAddressHexByte, decodeErr := base58.Decode(address)
		if decodeErr != nil {
			err = errors.WithMessage(err, fmt.Sprintf("decode address error%s", address))
		}
		tokenAddressHexStrRaw := hex.EncodeToString(tokenAddressHexByte)
		ret = fmt.Sprintf("0x%s", tokenAddressHexStrRaw)
		return
	}
	if evmType == "solana" {
		tokenAddressHexByte, decodeErr := base58.Decode(address)
		if decodeErr != nil {
			err = errors.WithMessage(err, fmt.Sprintf("decode address error%s", address))
		}
		tokenAddressHexStrRaw := hex.EncodeToString(tokenAddressHexByte)
		ret = fmt.Sprintf("0x%s", tokenAddressHexStrRaw)
		return
	}
	ret = address
	if !strings.HasPrefix(ret, "0x") {
		err = errors.WithMessage(GetNoEmptyError(err), "address format incorrect")
	}
	return
}
func GetTokenAddress(address string, evmType string) (ret string, err error) {
	ret = ""
	if evmType == "near" {
		tokenAddressHexByte, decodeErr := base58.Decode(address)
		if decodeErr != nil {
			err = errors.WithMessage(err, fmt.Sprintf("decode address error%s", address))
		}
		tokenAddressHexStrRaw := hex.EncodeToString(tokenAddressHexByte)
		ret = fmt.Sprintf("0x%s", tokenAddressHexStrRaw)
		return
	}
	if evmType == "solana" {
		tokenAddressHexByte, decodeErr := base58.Decode(address)
		if decodeErr != nil {
			err = errors.WithMessage(err, fmt.Sprintf("decode address error%s", address))
		}
		tokenAddressHexStrRaw := hex.EncodeToString(tokenAddressHexByte)
		ret = fmt.Sprintf("0x%s", tokenAddressHexStrRaw)
		return
	}
	ret = address
	if !strings.HasPrefix(ret, "0x") {
		err = errors.WithMessage(GetNoEmptyError(err), "address format incorrect")
	}
	return
}

package utils

import (
	"math/big"

	"github.com/shopspring/decimal"
)

func ConvertWeiToFloat(weiStr string, decimals int) float64 {
	if weiStr == "" || weiStr == "0" {
		return 0
	}

	wei := new(big.Int)
	wei.SetString(weiStr, 10)

	weiDecimal := decimal.NewFromBigInt(wei, 0)
	divisor := decimal.New(1, int32(decimals))
	result, _ := weiDecimal.Div(divisor).Float64()

	return result
}

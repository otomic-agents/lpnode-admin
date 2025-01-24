// Code generated by goa v3.11.0, DO NOT EDIT.
//
// tokenManager HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	tokenmanager "admin-panel/gen/token_manager"
	"encoding/json"
	"fmt"

	goa "goa.design/goa/v3/pkg"
)

// BuildTokenCreatePayload builds the payload for the tokenManager tokenCreate
// endpoint from CLI flags.
func BuildTokenCreatePayload(tokenManagerTokenCreateBody string) (*tokenmanager.TokenItem, error) {
	var err error
	var body TokenCreateRequestBody
	{
		err = json.Unmarshal([]byte(tokenManagerTokenCreateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"_id\": \"Autem neque distinctio dolor ut dolorem.\",\n      \"address\": \"Corrupti explicabo.\",\n      \"chainId\": 506831215610962050,\n      \"chainType\": \"Autem et autem.\",\n      \"coinType\": \"coin\",\n      \"marketName\": \"Aliquid ab ut.\",\n      \"precision\": 15,\n      \"tokenId\": \"Earum dicta.\",\n      \"tokenName\": \"Dolorum excepturi sit vero aperiam sit voluptas.\"\n   }'")
		}
		if body.Precision < 6 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.precision", body.Precision, 6, true))
		}
		if body.Precision > 18 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.precision", body.Precision, 18, false))
		}
		if !(body.CoinType == "stable_coin" || body.CoinType == "coin") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.coinType", body.CoinType, []interface{}{"stable_coin", "coin"}))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &tokenmanager.TokenItem{
		ID:         body.ID,
		TokenID:    body.TokenID,
		ChainID:    body.ChainID,
		Address:    body.Address,
		TokenName:  body.TokenName,
		MarketName: body.MarketName,
		Precision:  body.Precision,
		ChainType:  body.ChainType,
		CoinType:   body.CoinType,
	}

	return v, nil
}

// BuildTokenDeletePayload builds the payload for the tokenManager tokenDelete
// endpoint from CLI flags.
func BuildTokenDeletePayload(tokenManagerTokenDeleteBody string) (*tokenmanager.DeleteTokenFilter, error) {
	var err error
	var body TokenDeleteRequestBody
	{
		err = json.Unmarshal([]byte(tokenManagerTokenDeleteBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"_id\": \"Et odio numquam voluptatem iusto.\"\n   }'")
		}
	}
	v := &tokenmanager.DeleteTokenFilter{
		ID: body.ID,
	}

	return v, nil
}

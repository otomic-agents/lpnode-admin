// Code generated by goa v3.11.0, DO NOT EDIT.
//
// accountDex HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	accountdex "admin-panel/gen/account_dex"
	"encoding/json"
	"fmt"
)

// BuildWalletInfoPayload builds the payload for the accountDex walletInfo
// endpoint from CLI flags.
func BuildWalletInfoPayload(accountDexWalletInfoBody string) (*accountdex.WalletInfoPayload, error) {
	var err error
	var body WalletInfoRequestBody
	{
		err = json.Unmarshal([]byte(accountDexWalletInfoBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"chainId\": 517246796943780996\n   }'")
		}
	}
	v := &accountdex.WalletInfoPayload{
		ChainID: body.ChainID,
	}

	return v, nil
}

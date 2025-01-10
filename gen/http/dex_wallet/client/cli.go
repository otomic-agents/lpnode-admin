// Code generated by goa v3.11.0, DO NOT EDIT.
//
// dexWallet HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	dexwallet "admin-panel/gen/dex_wallet"
	"encoding/json"
	"fmt"

	goa "goa.design/goa/v3/pkg"
)

// BuildCreateDexWalletPayload builds the payload for the dexWallet
// createDexWallet endpoint from CLI flags.
func BuildCreateDexWalletPayload(dexWalletCreateDexWalletBody string) (*dexwallet.WalletRow, error) {
	var err error
	var body CreateDexWalletRequestBody
	{
		err = json.Unmarshal([]byte(dexWalletCreateDexWalletBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"accountId\": \"Voluptas et.\",\n      \"address\": \"Excepturi qui explicabo sit expedita fugit.\",\n      \"chainId\": 7078845938374597277,\n      \"chainType\": \"Qui quisquam illum non aut quaerat.\",\n      \"id\": \"Et vel dolor aut adipisci cupiditate.\",\n      \"privateKey\": \"Corrupti veritatis.\",\n      \"signServiceEndpoint\": \"Inventore voluptas officiis sed voluptates recusandae.\",\n      \"storeId\": \"A nulla ipsa.\",\n      \"vaultHostType\": \"Reprehenderit aut.\",\n      \"vaultName\": \"Repellat laudantium iure impedit nesciunt ut rerum.\",\n      \"vaultSecertType\": \"Aut et iusto voluptatem debitis earum voluptatem.\",\n      \"walletName\": \"Sint exercitationem quas debitis.\",\n      \"walletType\": \"storeId\"\n   }'")
		}
		if !(body.WalletType == "privateKey" || body.WalletType == "storeId") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.walletType", body.WalletType, []interface{}{"privateKey", "storeId"}))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &dexwallet.WalletRow{
		ID:                  body.ID,
		WalletName:          body.WalletName,
		PrivateKey:          body.PrivateKey,
		Address:             body.Address,
		ChainType:           body.ChainType,
		AccountID:           body.AccountID,
		ChainID:             body.ChainID,
		StoreID:             body.StoreID,
		VaultHostType:       body.VaultHostType,
		VaultName:           body.VaultName,
		VaultSecertType:     body.VaultSecertType,
		SignServiceEndpoint: body.SignServiceEndpoint,
		WalletType:          body.WalletType,
	}

	return v, nil
}

// BuildDeleteDexWalletPayload builds the payload for the dexWallet
// deleteDexWallet endpoint from CLI flags.
func BuildDeleteDexWalletPayload(dexWalletDeleteDexWalletBody string) (*dexwallet.DeleteFilter, error) {
	var err error
	var body DeleteDexWalletRequestBody
	{
		err = json.Unmarshal([]byte(dexWalletDeleteDexWalletBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"id\": \"Aut rerum repellendus.\"\n   }'")
		}
	}
	v := &dexwallet.DeleteFilter{
		ID: body.ID,
	}

	return v, nil
}

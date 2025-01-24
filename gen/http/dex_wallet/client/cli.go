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
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"accountId\": \"Fugiat dolores asperiores velit.\",\n      \"address\": \"Rerum illum recusandae.\",\n      \"balance\": \"Doloremque nam dolorum sint.\",\n      \"chainId\": 6285722081445360316,\n      \"chainType\": \"In ratione labore molestiae.\",\n      \"id\": \"Ut rerum praesentium omnis.\",\n      \"privateKey\": \"Vel dolores ullam incidunt labore rem quibusdam.\",\n      \"signServiceEndpoint\": \"Magnam sed.\",\n      \"storeId\": \"Consequuntur quod amet.\",\n      \"vaultHostType\": \"Ducimus itaque.\",\n      \"vaultName\": \"Ut ipsa et.\",\n      \"vaultSecertType\": \"Ab occaecati dignissimos cupiditate nisi.\",\n      \"walletName\": \"Aut rerum repellendus.\",\n      \"walletType\": \"storeId\"\n   }'")
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
		Balance:             body.Balance,
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
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"id\": \"Error ad perferendis ut unde.\"\n   }'")
		}
	}
	v := &dexwallet.DeleteFilter{
		ID: body.ID,
	}

	return v, nil
}

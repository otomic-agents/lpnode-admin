// Code generated by goa v3.11.0, DO NOT EDIT.
//
// dexWallet HTTP server types
//
// Command:
// $ goa gen admin-panel/design

package server

import (
	dexwallet "admin-panel/gen/dex_wallet"

	goa "goa.design/goa/v3/pkg"
)

// CreateDexWalletRequestBody is the type of the "dexWallet" service
// "createDexWallet" endpoint HTTP request body.
type CreateDexWalletRequestBody struct {
	// mongodb primary key
	ID         *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	WalletName *string `form:"walletName,omitempty" json:"walletName,omitempty" xml:"walletName,omitempty"`
	PrivateKey *string `form:"privateKey,omitempty" json:"privateKey,omitempty" xml:"privateKey,omitempty"`
	Address    *string `form:"address,omitempty" json:"address,omitempty" xml:"address,omitempty"`
	ChainType  *string `form:"chainType,omitempty" json:"chainType,omitempty" xml:"chainType,omitempty"`
	// wallet
	AccountID *string `form:"accountId,omitempty" json:"accountId,omitempty" xml:"accountId,omitempty"`
	// chain Id
	ChainID             *int64  `form:"chainId,omitempty" json:"chainId,omitempty" xml:"chainId,omitempty"`
	StoreID             *string `form:"storeId,omitempty" json:"storeId,omitempty" xml:"storeId,omitempty"`
	VaultHostType       *string `form:"vaultHostType,omitempty" json:"vaultHostType,omitempty" xml:"vaultHostType,omitempty"`
	VaultName           *string `form:"vaultName,omitempty" json:"vaultName,omitempty" xml:"vaultName,omitempty"`
	VaultSecertType     *string `form:"vaultSecertType,omitempty" json:"vaultSecertType,omitempty" xml:"vaultSecertType,omitempty"`
	SignServiceEndpoint *string `form:"signServiceEndpoint,omitempty" json:"signServiceEndpoint,omitempty" xml:"signServiceEndpoint,omitempty"`
	WalletType          *string `form:"walletType,omitempty" json:"walletType,omitempty" xml:"walletType,omitempty"`
	Balance             *string `form:"balance,omitempty" json:"balance,omitempty" xml:"balance,omitempty"`
}

// DeleteDexWalletRequestBody is the type of the "dexWallet" service
// "deleteDexWallet" endpoint HTTP request body.
type DeleteDexWalletRequestBody struct {
	// mongodb primary key
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
}

// ListDexWalletResponseBody is the type of the "dexWallet" service
// "listDexWallet" endpoint HTTP response body.
type ListDexWalletResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// wallet list
	Result  []*WalletRowResponseBody `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string                  `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// CreateDexWalletResponseBody is the type of the "dexWallet" service
// "createDexWallet" endpoint HTTP response body.
type CreateDexWalletResponseBody struct {
	Code   *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	Result *struct {
		ID *string `form:"_id" json:"_id" xml:"_id"`
	} `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// DeleteDexWalletResponseBody is the type of the "dexWallet" service
// "deleteDexWallet" endpoint HTTP response body.
type DeleteDexWalletResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// result
	Result  *int64  `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// VaultListResponseBody is the type of the "dexWallet" service "vaultList"
// endpoint HTTP response body.
type VaultListResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// list
	Result  []*VaultRowResponseBody `form:"result" json:"result" xml:"result"`
	Message *string                 `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// UpdateLpWalletResponseBody is the type of the "dexWallet" service
// "updateLpWallet" endpoint HTTP response body.
type UpdateLpWalletResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// list
	Result  string  `form:"result" json:"result" xml:"result"`
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// WalletRowResponseBody is used to define fields on response body types.
type WalletRowResponseBody struct {
	// mongodb primary key
	ID         *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	WalletName string  `form:"walletName" json:"walletName" xml:"walletName"`
	PrivateKey *string `form:"privateKey,omitempty" json:"privateKey,omitempty" xml:"privateKey,omitempty"`
	Address    *string `form:"address,omitempty" json:"address,omitempty" xml:"address,omitempty"`
	ChainType  string  `form:"chainType" json:"chainType" xml:"chainType"`
	// wallet
	AccountID *string `form:"accountId,omitempty" json:"accountId,omitempty" xml:"accountId,omitempty"`
	// chain Id
	ChainID             int64   `form:"chainId" json:"chainId" xml:"chainId"`
	StoreID             *string `form:"storeId,omitempty" json:"storeId,omitempty" xml:"storeId,omitempty"`
	VaultHostType       *string `form:"vaultHostType,omitempty" json:"vaultHostType,omitempty" xml:"vaultHostType,omitempty"`
	VaultName           *string `form:"vaultName,omitempty" json:"vaultName,omitempty" xml:"vaultName,omitempty"`
	VaultSecertType     *string `form:"vaultSecertType,omitempty" json:"vaultSecertType,omitempty" xml:"vaultSecertType,omitempty"`
	SignServiceEndpoint *string `form:"signServiceEndpoint,omitempty" json:"signServiceEndpoint,omitempty" xml:"signServiceEndpoint,omitempty"`
	WalletType          string  `form:"walletType" json:"walletType" xml:"walletType"`
	Balance             *string `form:"balance,omitempty" json:"balance,omitempty" xml:"balance,omitempty"`
}

// VaultRowResponseBody is used to define fields on response body types.
type VaultRowResponseBody struct {
	// address
	Address *string `form:"address,omitempty" json:"address,omitempty" xml:"address,omitempty"`
	// host type
	HostType *string `form:"hostType,omitempty" json:"hostType,omitempty" xml:"hostType,omitempty"`
	// storeId
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// store name
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// store secert type
	SecertType *string `form:"secertType,omitempty" json:"secertType,omitempty" xml:"secertType,omitempty"`
}

// NewListDexWalletResponseBody builds the HTTP response body from the result
// of the "listDexWallet" endpoint of the "dexWallet" service.
func NewListDexWalletResponseBody(res *dexwallet.ListDexWalletResult) *ListDexWalletResponseBody {
	body := &ListDexWalletResponseBody{
		Code:    res.Code,
		Message: res.Message,
	}
	if res.Result != nil {
		body.Result = make([]*WalletRowResponseBody, len(res.Result))
		for i, val := range res.Result {
			body.Result[i] = marshalDexwalletWalletRowToWalletRowResponseBody(val)
		}
	}
	return body
}

// NewCreateDexWalletResponseBody builds the HTTP response body from the result
// of the "createDexWallet" endpoint of the "dexWallet" service.
func NewCreateDexWalletResponseBody(res *dexwallet.CreateDexWalletResult) *CreateDexWalletResponseBody {
	body := &CreateDexWalletResponseBody{
		Code:    res.Code,
		Message: res.Message,
	}
	if res.Result != nil {
		body.Result = &struct {
			ID *string `form:"_id" json:"_id" xml:"_id"`
		}{
			ID: res.Result.ID,
		}
	}
	return body
}

// NewDeleteDexWalletResponseBody builds the HTTP response body from the result
// of the "deleteDexWallet" endpoint of the "dexWallet" service.
func NewDeleteDexWalletResponseBody(res *dexwallet.DeleteDexWalletResult) *DeleteDexWalletResponseBody {
	body := &DeleteDexWalletResponseBody{
		Code:    res.Code,
		Result:  res.Result,
		Message: res.Message,
	}
	return body
}

// NewVaultListResponseBody builds the HTTP response body from the result of
// the "vaultList" endpoint of the "dexWallet" service.
func NewVaultListResponseBody(res *dexwallet.VaultListResult) *VaultListResponseBody {
	body := &VaultListResponseBody{
		Code:    res.Code,
		Message: res.Message,
	}
	if res.Result != nil {
		body.Result = make([]*VaultRowResponseBody, len(res.Result))
		for i, val := range res.Result {
			body.Result[i] = marshalDexwalletVaultRowToVaultRowResponseBody(val)
		}
	}
	return body
}

// NewUpdateLpWalletResponseBody builds the HTTP response body from the result
// of the "updateLpWallet" endpoint of the "dexWallet" service.
func NewUpdateLpWalletResponseBody(res *dexwallet.UpdateLpWalletResult) *UpdateLpWalletResponseBody {
	body := &UpdateLpWalletResponseBody{
		Code:    res.Code,
		Result:  res.Result,
		Message: res.Message,
	}
	return body
}

// NewCreateDexWalletWalletRow builds a dexWallet service createDexWallet
// endpoint payload.
func NewCreateDexWalletWalletRow(body *CreateDexWalletRequestBody) *dexwallet.WalletRow {
	v := &dexwallet.WalletRow{
		ID:                  body.ID,
		WalletName:          *body.WalletName,
		PrivateKey:          body.PrivateKey,
		Address:             body.Address,
		ChainType:           *body.ChainType,
		AccountID:           body.AccountID,
		ChainID:             *body.ChainID,
		StoreID:             body.StoreID,
		VaultHostType:       body.VaultHostType,
		VaultName:           body.VaultName,
		VaultSecertType:     body.VaultSecertType,
		SignServiceEndpoint: body.SignServiceEndpoint,
		WalletType:          *body.WalletType,
		Balance:             body.Balance,
	}

	return v
}

// NewDeleteDexWalletDeleteFilter builds a dexWallet service deleteDexWallet
// endpoint payload.
func NewDeleteDexWalletDeleteFilter(body *DeleteDexWalletRequestBody) *dexwallet.DeleteFilter {
	v := &dexwallet.DeleteFilter{
		ID: *body.ID,
	}

	return v
}

// ValidateCreateDexWalletRequestBody runs the validations defined on
// CreateDexWalletRequestBody
func ValidateCreateDexWalletRequestBody(body *CreateDexWalletRequestBody) (err error) {
	if body.WalletName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("walletName", "body"))
	}
	if body.ChainID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("chainId", "body"))
	}
	if body.ChainType == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("chainType", "body"))
	}
	if body.WalletType == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("walletType", "body"))
	}
	if body.WalletType != nil {
		if !(*body.WalletType == "privateKey" || *body.WalletType == "storeId") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.walletType", *body.WalletType, []interface{}{"privateKey", "storeId"}))
		}
	}
	return
}

// ValidateDeleteDexWalletRequestBody runs the validations defined on
// DeleteDexWalletRequestBody
func ValidateDeleteDexWalletRequestBody(body *DeleteDexWalletRequestBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	return
}

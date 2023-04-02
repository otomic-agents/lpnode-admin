// Code generated by goa v3.11.0, DO NOT EDIT.
//
// dexWallet HTTP client types
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	dexwallet "admin-panel/gen/dex_wallet"

	goa "goa.design/goa/v3/pkg"
)

// CreateDexWalletRequestBody is the type of the "dexWallet" service
// "createDexWallet" endpoint HTTP request body.
type CreateDexWalletRequestBody struct {
	// mongodb主键
	ID         *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	WalletName string  `form:"walletName" json:"walletName" xml:"walletName"`
	PrivateKey *string `form:"privateKey,omitempty" json:"privateKey,omitempty" xml:"privateKey,omitempty"`
	Address    *string `form:"address,omitempty" json:"address,omitempty" xml:"address,omitempty"`
	ChainType  string  `form:"chainType" json:"chainType" xml:"chainType"`
	// wallet对应的人类可阅读的名称
	AccountID *string `form:"accountId,omitempty" json:"accountId,omitempty" xml:"accountId,omitempty"`
	// 链的Id
	ChainID         int64   `form:"chainId" json:"chainId" xml:"chainId"`
	StoreID         *string `form:"storeId,omitempty" json:"storeId,omitempty" xml:"storeId,omitempty"`
	VaultHostType   *string `form:"vaultHostType,omitempty" json:"vaultHostType,omitempty" xml:"vaultHostType,omitempty"`
	VaultName       *string `form:"vaultName,omitempty" json:"vaultName,omitempty" xml:"vaultName,omitempty"`
	VaultSecertType *string `form:"vaultSecertType,omitempty" json:"vaultSecertType,omitempty" xml:"vaultSecertType,omitempty"`
	WalletType      string  `form:"walletType" json:"walletType" xml:"walletType"`
}

// DeleteDexWalletRequestBody is the type of the "dexWallet" service
// "deleteDexWallet" endpoint HTTP request body.
type DeleteDexWalletRequestBody struct {
	// Mongodb 的主键
	ID string `form:"id" json:"id" xml:"id"`
}

// ListDexWalletResponseBody is the type of the "dexWallet" service
// "listDexWallet" endpoint HTTP response body.
type ListDexWalletResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// 钱包别表
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
	// 是否删除成功
	Result  *int64  `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// VaultListResponseBody is the type of the "dexWallet" service "vaultList"
// endpoint HTTP response body.
type VaultListResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// 列表
	Result  []*VaultRowResponseBody `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string                 `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// WalletRowResponseBody is used to define fields on response body types.
type WalletRowResponseBody struct {
	// mongodb主键
	ID         *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	WalletName *string `form:"walletName,omitempty" json:"walletName,omitempty" xml:"walletName,omitempty"`
	PrivateKey *string `form:"privateKey,omitempty" json:"privateKey,omitempty" xml:"privateKey,omitempty"`
	Address    *string `form:"address,omitempty" json:"address,omitempty" xml:"address,omitempty"`
	ChainType  *string `form:"chainType,omitempty" json:"chainType,omitempty" xml:"chainType,omitempty"`
	// wallet对应的人类可阅读的名称
	AccountID *string `form:"accountId,omitempty" json:"accountId,omitempty" xml:"accountId,omitempty"`
	// 链的Id
	ChainID         *int64  `form:"chainId,omitempty" json:"chainId,omitempty" xml:"chainId,omitempty"`
	StoreID         *string `form:"storeId,omitempty" json:"storeId,omitempty" xml:"storeId,omitempty"`
	VaultHostType   *string `form:"vaultHostType,omitempty" json:"vaultHostType,omitempty" xml:"vaultHostType,omitempty"`
	VaultName       *string `form:"vaultName,omitempty" json:"vaultName,omitempty" xml:"vaultName,omitempty"`
	VaultSecertType *string `form:"vaultSecertType,omitempty" json:"vaultSecertType,omitempty" xml:"vaultSecertType,omitempty"`
	WalletType      *string `form:"walletType,omitempty" json:"walletType,omitempty" xml:"walletType,omitempty"`
}

// VaultRowResponseBody is used to define fields on response body types.
type VaultRowResponseBody struct {
	// 地址
	Address *string `form:"address,omitempty" json:"address,omitempty" xml:"address,omitempty"`
	// 托管类型
	HostType *string `form:"hostType,omitempty" json:"hostType,omitempty" xml:"hostType,omitempty"`
	// 存储Id
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// 钱包名称
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// 私钥类型
	SecertType *string `form:"secertType,omitempty" json:"secertType,omitempty" xml:"secertType,omitempty"`
}

// NewCreateDexWalletRequestBody builds the HTTP request body from the payload
// of the "createDexWallet" endpoint of the "dexWallet" service.
func NewCreateDexWalletRequestBody(p *dexwallet.WalletRow) *CreateDexWalletRequestBody {
	body := &CreateDexWalletRequestBody{
		ID:              p.ID,
		WalletName:      p.WalletName,
		PrivateKey:      p.PrivateKey,
		Address:         p.Address,
		ChainType:       p.ChainType,
		AccountID:       p.AccountID,
		ChainID:         p.ChainID,
		StoreID:         p.StoreID,
		VaultHostType:   p.VaultHostType,
		VaultName:       p.VaultName,
		VaultSecertType: p.VaultSecertType,
		WalletType:      p.WalletType,
	}
	return body
}

// NewDeleteDexWalletRequestBody builds the HTTP request body from the payload
// of the "deleteDexWallet" endpoint of the "dexWallet" service.
func NewDeleteDexWalletRequestBody(p *dexwallet.DeleteFilter) *DeleteDexWalletRequestBody {
	body := &DeleteDexWalletRequestBody{
		ID: p.ID,
	}
	return body
}

// NewListDexWalletResultOK builds a "dexWallet" service "listDexWallet"
// endpoint result from a HTTP "OK" response.
func NewListDexWalletResultOK(body *ListDexWalletResponseBody) *dexwallet.ListDexWalletResult {
	v := &dexwallet.ListDexWalletResult{
		Code:    body.Code,
		Message: body.Message,
	}
	if body.Result != nil {
		v.Result = make([]*dexwallet.WalletRow, len(body.Result))
		for i, val := range body.Result {
			v.Result[i] = unmarshalWalletRowResponseBodyToDexwalletWalletRow(val)
		}
	}

	return v
}

// NewCreateDexWalletResultOK builds a "dexWallet" service "createDexWallet"
// endpoint result from a HTTP "OK" response.
func NewCreateDexWalletResultOK(body *CreateDexWalletResponseBody) *dexwallet.CreateDexWalletResult {
	v := &dexwallet.CreateDexWalletResult{
		Code:    body.Code,
		Message: body.Message,
	}
	if body.Result != nil {
		v.Result = &struct {
			ID *string
		}{
			ID: body.Result.ID,
		}
	}

	return v
}

// NewDeleteDexWalletResultOK builds a "dexWallet" service "deleteDexWallet"
// endpoint result from a HTTP "OK" response.
func NewDeleteDexWalletResultOK(body *DeleteDexWalletResponseBody) *dexwallet.DeleteDexWalletResult {
	v := &dexwallet.DeleteDexWalletResult{
		Code:    body.Code,
		Result:  body.Result,
		Message: body.Message,
	}

	return v
}

// NewVaultListResultOK builds a "dexWallet" service "vaultList" endpoint
// result from a HTTP "OK" response.
func NewVaultListResultOK(body *VaultListResponseBody) *dexwallet.VaultListResult {
	v := &dexwallet.VaultListResult{
		Code:    body.Code,
		Message: body.Message,
	}
	if body.Result != nil {
		v.Result = make([]*dexwallet.VaultRow, len(body.Result))
		for i, val := range body.Result {
			v.Result[i] = unmarshalVaultRowResponseBodyToDexwalletVaultRow(val)
		}
	}

	return v
}

// ValidateListDexWalletResponseBody runs the validations defined on
// ListDexWalletResponseBody
func ValidateListDexWalletResponseBody(body *ListDexWalletResponseBody) (err error) {
	for _, e := range body.Result {
		if e != nil {
			if err2 := ValidateWalletRowResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateWalletRowResponseBody runs the validations defined on
// walletRowResponseBody
func ValidateWalletRowResponseBody(body *WalletRowResponseBody) (err error) {
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

// Code generated by goa v3.11.0, DO NOT EDIT.
//
// bridgeConfig HTTP server types
//
// Command:
// $ goa gen admin-panel/design

package server

import (
	bridgeconfig "admin-panel/gen/bridge_config"

	goa "goa.design/goa/v3/pkg"
)

// BridgeCreateRequestBody is the type of the "bridgeConfig" service
// "bridgeCreate" endpoint HTTP request body.
type BridgeCreateRequestBody struct {
	// bridge name ****
	BridgeName *string `form:"bridgeName,omitempty" json:"bridgeName,omitempty" xml:"bridgeName,omitempty"`
	// mongodb primary key, from basedata
	SrcChainID *string `form:"srcChainId,omitempty" json:"srcChainId,omitempty" xml:"srcChainId,omitempty"`
	// mongodb primary key, from basedata
	DstChainID *string `form:"dstChainId,omitempty" json:"dstChainId,omitempty" xml:"dstChainId,omitempty"`
	// mongodb primary key, from tokenlist
	SrcTokenID *string `form:"srcTokenId,omitempty" json:"srcTokenId,omitempty" xml:"srcTokenId,omitempty"`
	// mongodb primary key, from tokenlist
	DstTokenID *string `form:"dstTokenId,omitempty" json:"dstTokenId,omitempty" xml:"dstTokenId,omitempty"`
	// mongodb primary key, from walletlist
	WalletID *string `form:"walletId,omitempty" json:"walletId,omitempty" xml:"walletId,omitempty"`
	// mongodb primary key, from walletlist
	SrcWalletID *string `form:"srcWalletId,omitempty" json:"srcWalletId,omitempty" xml:"srcWalletId,omitempty"`
	// amm name at install
	AmmName *string `form:"ammName,omitempty" json:"ammName,omitempty" xml:"ammName,omitempty"`
	// relay api key
	RelayAPIKey *string `form:"relayApiKey,omitempty" json:"relayApiKey,omitempty" xml:"relayApiKey,omitempty"`
	// relayUri
	RelayURI      *string `form:"relayUri,omitempty" json:"relayUri,omitempty" xml:"relayUri,omitempty"`
	EnableHedge   *bool   `form:"enableHedge,omitempty" json:"enableHedge,omitempty" xml:"enableHedge,omitempty"`
	EnableLimiter *bool   `form:"enableLimiter,omitempty" json:"enableLimiter,omitempty" xml:"enableLimiter,omitempty"`
}

// BridgeDeleteRequestBody is the type of the "bridgeConfig" service
// "bridgeDelete" endpoint HTTP request body.
type BridgeDeleteRequestBody struct {
	// mongodb primary key
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
}

// BridgeTestRequestBody is the type of the "bridgeConfig" service "bridgeTest"
// endpoint HTTP request body.
type BridgeTestRequestBody struct {
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
}

// BridgeCreateResponseBody is the type of the "bridgeConfig" service
// "bridgeCreate" endpoint HTTP response body.
type BridgeCreateResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// result
	Result  *int64  `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// BridgeListResponseBody is the type of the "bridgeConfig" service
// "bridgeList" endpoint HTTP response body.
type BridgeListResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// chain list
	Result  []*ListBridgeItemResponseBody `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string                       `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// BridgeDeleteResponseBody is the type of the "bridgeConfig" service
// "bridgeDelete" endpoint HTTP response body.
type BridgeDeleteResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// result
	Result  *int64  `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// BridgeTestResponseBody is the type of the "bridgeConfig" service
// "bridgeTest" endpoint HTTP response body.
type BridgeTestResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
}

// ListBridgeItemResponseBody is used to define fields on response body types.
type ListBridgeItemResponseBody struct {
	ID                *string `form:"_id,omitempty" json:"_id,omitempty" xml:"_id,omitempty"`
	DstChainID        *string `form:"dstChainId,omitempty" json:"dstChainId,omitempty" xml:"dstChainId,omitempty"`
	DstTokenID        *string `form:"dstTokenId,omitempty" json:"dstTokenId,omitempty" xml:"dstTokenId,omitempty"`
	SrcChainID        *string `form:"srcChainId,omitempty" json:"srcChainId,omitempty" xml:"srcChainId,omitempty"`
	SrcTokenID        *string `form:"srcTokenId,omitempty" json:"srcTokenId,omitempty" xml:"srcTokenId,omitempty"`
	AmmName           *string `form:"ammName,omitempty" json:"ammName,omitempty" xml:"ammName,omitempty"`
	BridgeName        *string `form:"bridgeName,omitempty" json:"bridgeName,omitempty" xml:"bridgeName,omitempty"`
	DstChainRawID     *int64  `form:"dstChainRawId,omitempty" json:"dstChainRawId,omitempty" xml:"dstChainRawId,omitempty"`
	DstClientURI      *string `form:"dstClientUri,omitempty" json:"dstClientUri,omitempty" xml:"dstClientUri,omitempty"`
	DstToken          *string `form:"dstToken,omitempty" json:"dstToken,omitempty" xml:"dstToken,omitempty"`
	LpReceiverAddress *string `form:"lpReceiverAddress,omitempty" json:"lpReceiverAddress,omitempty" xml:"lpReceiverAddress,omitempty"`
	MsmqName          *string `form:"msmqName,omitempty" json:"msmqName,omitempty" xml:"msmqName,omitempty"`
	SrcChainRawID     *int64  `form:"srcChainRawId,omitempty" json:"srcChainRawId,omitempty" xml:"srcChainRawId,omitempty"`
	SrcToken          *string `form:"srcToken,omitempty" json:"srcToken,omitempty" xml:"srcToken,omitempty"`
	WalletName        *string `form:"walletName,omitempty" json:"walletName,omitempty" xml:"walletName,omitempty"`
	WalletID          *string `form:"walletId,omitempty" json:"walletId,omitempty" xml:"walletId,omitempty"`
	EnableHedge       *bool   `form:"enableHedge,omitempty" json:"enableHedge,omitempty" xml:"enableHedge,omitempty"`
	// Source chain token balance
	SrcTokenBalance string `form:"srcTokenBalance" json:"srcTokenBalance" xml:"srcTokenBalance"`
	// Destination chain token balance
	DstTokenBalance string `form:"dstTokenBalance" json:"dstTokenBalance" xml:"dstTokenBalance"`
	// Source token decimals
	SrcTokenDecimals int64 `form:"srcTokenDecimals" json:"srcTokenDecimals" xml:"srcTokenDecimals"`
	// Destination token decimals
	DstTokenDecimals int64 `form:"dstTokenDecimals" json:"dstTokenDecimals" xml:"dstTokenDecimals"`
}

// NewBridgeCreateResponseBody builds the HTTP response body from the result of
// the "bridgeCreate" endpoint of the "bridgeConfig" service.
func NewBridgeCreateResponseBody(res *bridgeconfig.BridgeCreateResult) *BridgeCreateResponseBody {
	body := &BridgeCreateResponseBody{
		Code:    res.Code,
		Result:  res.Result,
		Message: res.Message,
	}
	return body
}

// NewBridgeListResponseBody builds the HTTP response body from the result of
// the "bridgeList" endpoint of the "bridgeConfig" service.
func NewBridgeListResponseBody(res *bridgeconfig.BridgeListResult) *BridgeListResponseBody {
	body := &BridgeListResponseBody{
		Code:    res.Code,
		Message: res.Message,
	}
	if res.Result != nil {
		body.Result = make([]*ListBridgeItemResponseBody, len(res.Result))
		for i, val := range res.Result {
			body.Result[i] = marshalBridgeconfigListBridgeItemToListBridgeItemResponseBody(val)
		}
	}
	return body
}

// NewBridgeDeleteResponseBody builds the HTTP response body from the result of
// the "bridgeDelete" endpoint of the "bridgeConfig" service.
func NewBridgeDeleteResponseBody(res *bridgeconfig.BridgeDeleteResult) *BridgeDeleteResponseBody {
	body := &BridgeDeleteResponseBody{
		Code:    res.Code,
		Result:  res.Result,
		Message: res.Message,
	}
	return body
}

// NewBridgeTestResponseBody builds the HTTP response body from the result of
// the "bridgeTest" endpoint of the "bridgeConfig" service.
func NewBridgeTestResponseBody(res *bridgeconfig.BridgeTestResult) *BridgeTestResponseBody {
	body := &BridgeTestResponseBody{
		Code: res.Code,
	}
	return body
}

// NewBridgeCreateBridgeItem builds a bridgeConfig service bridgeCreate
// endpoint payload.
func NewBridgeCreateBridgeItem(body *BridgeCreateRequestBody) *bridgeconfig.BridgeItem {
	v := &bridgeconfig.BridgeItem{
		BridgeName:  *body.BridgeName,
		SrcChainID:  *body.SrcChainID,
		DstChainID:  *body.DstChainID,
		SrcTokenID:  *body.SrcTokenID,
		DstTokenID:  *body.DstTokenID,
		WalletID:    *body.WalletID,
		SrcWalletID: *body.SrcWalletID,
		AmmName:     *body.AmmName,
		RelayAPIKey: *body.RelayAPIKey,
		RelayURI:    *body.RelayURI,
	}
	if body.EnableHedge != nil {
		v.EnableHedge = *body.EnableHedge
	}
	if body.EnableLimiter != nil {
		v.EnableLimiter = *body.EnableLimiter
	}
	if body.EnableHedge == nil {
		v.EnableHedge = true
	}
	if body.EnableLimiter == nil {
		v.EnableLimiter = true
	}

	return v
}

// NewBridgeDeleteDeleteBridgeFilter builds a bridgeConfig service bridgeDelete
// endpoint payload.
func NewBridgeDeleteDeleteBridgeFilter(body *BridgeDeleteRequestBody) *bridgeconfig.DeleteBridgeFilter {
	v := &bridgeconfig.DeleteBridgeFilter{
		ID: *body.ID,
	}

	return v
}

// NewBridgeTestPayload builds a bridgeConfig service bridgeTest endpoint
// payload.
func NewBridgeTestPayload(body *BridgeTestRequestBody) *bridgeconfig.BridgeTestPayload {
	v := &bridgeconfig.BridgeTestPayload{
		ID: body.ID,
	}

	return v
}

// ValidateBridgeCreateRequestBody runs the validations defined on
// BridgeCreateRequestBody
func ValidateBridgeCreateRequestBody(body *BridgeCreateRequestBody) (err error) {
	if body.BridgeName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bridgeName", "body"))
	}
	if body.SrcChainID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("srcChainId", "body"))
	}
	if body.DstChainID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("dstChainId", "body"))
	}
	if body.SrcTokenID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("srcTokenId", "body"))
	}
	if body.DstTokenID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("dstTokenId", "body"))
	}
	if body.SrcWalletID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("srcWalletId", "body"))
	}
	if body.WalletID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("walletId", "body"))
	}
	if body.AmmName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("ammName", "body"))
	}
	if body.RelayAPIKey == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("relayApiKey", "body"))
	}
	if body.RelayURI == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("relayUri", "body"))
	}
	return
}

// ValidateBridgeDeleteRequestBody runs the validations defined on
// BridgeDeleteRequestBody
func ValidateBridgeDeleteRequestBody(body *BridgeDeleteRequestBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	return
}

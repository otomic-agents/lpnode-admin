// Code generated by goa v3.11.0, DO NOT EDIT.
//
// dexWallet client HTTP transport
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the dexWallet service endpoint HTTP clients.
type Client struct {
	// ListDexWallet Doer is the HTTP client used to make requests to the
	// listDexWallet endpoint.
	ListDexWalletDoer goahttp.Doer

	// CreateDexWallet Doer is the HTTP client used to make requests to the
	// createDexWallet endpoint.
	CreateDexWalletDoer goahttp.Doer

	// DeleteDexWallet Doer is the HTTP client used to make requests to the
	// deleteDexWallet endpoint.
	DeleteDexWalletDoer goahttp.Doer

	// VaultList Doer is the HTTP client used to make requests to the vaultList
	// endpoint.
	VaultListDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the dexWallet service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		ListDexWalletDoer:   doer,
		CreateDexWalletDoer: doer,
		DeleteDexWalletDoer: doer,
		VaultListDoer:       doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// ListDexWallet returns an endpoint that makes HTTP requests to the dexWallet
// service listDexWallet server.
func (c *Client) ListDexWallet() goa.Endpoint {
	var (
		decodeResponse = DecodeListDexWalletResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildListDexWalletRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ListDexWalletDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("dexWallet", "listDexWallet", err)
		}
		return decodeResponse(resp)
	}
}

// CreateDexWallet returns an endpoint that makes HTTP requests to the
// dexWallet service createDexWallet server.
func (c *Client) CreateDexWallet() goa.Endpoint {
	var (
		encodeRequest  = EncodeCreateDexWalletRequest(c.encoder)
		decodeResponse = DecodeCreateDexWalletResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildCreateDexWalletRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.CreateDexWalletDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("dexWallet", "createDexWallet", err)
		}
		return decodeResponse(resp)
	}
}

// DeleteDexWallet returns an endpoint that makes HTTP requests to the
// dexWallet service deleteDexWallet server.
func (c *Client) DeleteDexWallet() goa.Endpoint {
	var (
		encodeRequest  = EncodeDeleteDexWalletRequest(c.encoder)
		decodeResponse = DecodeDeleteDexWalletResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDeleteDexWalletRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DeleteDexWalletDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("dexWallet", "deleteDexWallet", err)
		}
		return decodeResponse(resp)
	}
}

// VaultList returns an endpoint that makes HTTP requests to the dexWallet
// service vaultList server.
func (c *Client) VaultList() goa.Endpoint {
	var (
		decodeResponse = DecodeVaultListResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildVaultListRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.VaultListDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("dexWallet", "vaultList", err)
		}
		return decodeResponse(resp)
	}
}

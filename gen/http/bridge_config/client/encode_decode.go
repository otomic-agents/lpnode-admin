// Code generated by goa v3.11.0, DO NOT EDIT.
//
// bridgeConfig HTTP client encoders and decoders
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	bridgeconfig "admin-panel/gen/bridge_config"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/v3/http"
)

// BuildBridgeCreateRequest instantiates a HTTP request object with method and
// path set to call the "bridgeConfig" service "bridgeCreate" endpoint
func (c *Client) BuildBridgeCreateRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: BridgeCreateBridgeConfigPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("bridgeConfig", "bridgeCreate", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeBridgeCreateRequest returns an encoder for requests sent to the
// bridgeConfig bridgeCreate server.
func EncodeBridgeCreateRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*bridgeconfig.BridgeItem)
		if !ok {
			return goahttp.ErrInvalidType("bridgeConfig", "bridgeCreate", "*bridgeconfig.BridgeItem", v)
		}
		body := NewBridgeCreateRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("bridgeConfig", "bridgeCreate", err)
		}
		return nil
	}
}

// DecodeBridgeCreateResponse returns a decoder for responses returned by the
// bridgeConfig bridgeCreate endpoint. restoreBody controls whether the
// response body should be restored after having been read.
func DecodeBridgeCreateResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body BridgeCreateResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("bridgeConfig", "bridgeCreate", err)
			}
			res := NewBridgeCreateResultOK(&body)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("bridgeConfig", "bridgeCreate", resp.StatusCode, string(body))
		}
	}
}

// BuildBridgeListRequest instantiates a HTTP request object with method and
// path set to call the "bridgeConfig" service "bridgeList" endpoint
func (c *Client) BuildBridgeListRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: BridgeListBridgeConfigPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("bridgeConfig", "bridgeList", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeBridgeListResponse returns a decoder for responses returned by the
// bridgeConfig bridgeList endpoint. restoreBody controls whether the response
// body should be restored after having been read.
func DecodeBridgeListResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body BridgeListResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("bridgeConfig", "bridgeList", err)
			}
			res := NewBridgeListResultOK(&body)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("bridgeConfig", "bridgeList", resp.StatusCode, string(body))
		}
	}
}

// BuildBridgeDeleteRequest instantiates a HTTP request object with method and
// path set to call the "bridgeConfig" service "bridgeDelete" endpoint
func (c *Client) BuildBridgeDeleteRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: BridgeDeleteBridgeConfigPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("bridgeConfig", "bridgeDelete", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeBridgeDeleteRequest returns an encoder for requests sent to the
// bridgeConfig bridgeDelete server.
func EncodeBridgeDeleteRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*bridgeconfig.DeleteBridgeFilter)
		if !ok {
			return goahttp.ErrInvalidType("bridgeConfig", "bridgeDelete", "*bridgeconfig.DeleteBridgeFilter", v)
		}
		body := NewBridgeDeleteRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("bridgeConfig", "bridgeDelete", err)
		}
		return nil
	}
}

// DecodeBridgeDeleteResponse returns a decoder for responses returned by the
// bridgeConfig bridgeDelete endpoint. restoreBody controls whether the
// response body should be restored after having been read.
func DecodeBridgeDeleteResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body BridgeDeleteResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("bridgeConfig", "bridgeDelete", err)
			}
			res := NewBridgeDeleteResultOK(&body)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("bridgeConfig", "bridgeDelete", resp.StatusCode, string(body))
		}
	}
}

// BuildBridgeTestRequest instantiates a HTTP request object with method and
// path set to call the "bridgeConfig" service "bridgeTest" endpoint
func (c *Client) BuildBridgeTestRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: BridgeTestBridgeConfigPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("bridgeConfig", "bridgeTest", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeBridgeTestRequest returns an encoder for requests sent to the
// bridgeConfig bridgeTest server.
func EncodeBridgeTestRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*bridgeconfig.BridgeTestPayload)
		if !ok {
			return goahttp.ErrInvalidType("bridgeConfig", "bridgeTest", "*bridgeconfig.BridgeTestPayload", v)
		}
		body := NewBridgeTestRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("bridgeConfig", "bridgeTest", err)
		}
		return nil
	}
}

// DecodeBridgeTestResponse returns a decoder for responses returned by the
// bridgeConfig bridgeTest endpoint. restoreBody controls whether the response
// body should be restored after having been read.
func DecodeBridgeTestResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body BridgeTestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("bridgeConfig", "bridgeTest", err)
			}
			res := NewBridgeTestResultOK(&body)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("bridgeConfig", "bridgeTest", resp.StatusCode, string(body))
		}
	}
}

// unmarshalListBridgeItemResponseBodyToBridgeconfigListBridgeItem builds a
// value of type *bridgeconfig.ListBridgeItem from a value of type
// *ListBridgeItemResponseBody.
func unmarshalListBridgeItemResponseBodyToBridgeconfigListBridgeItem(v *ListBridgeItemResponseBody) *bridgeconfig.ListBridgeItem {
	if v == nil {
		return nil
	}
	res := &bridgeconfig.ListBridgeItem{
		ID:                v.ID,
		DstChainID:        v.DstChainID,
		DstTokenID:        v.DstTokenID,
		SrcChainID:        v.SrcChainID,
		SrcTokenID:        v.SrcTokenID,
		AmmName:           v.AmmName,
		BridgeName:        v.BridgeName,
		DstChainRawID:     v.DstChainRawID,
		DstClientURI:      v.DstClientURI,
		DstToken:          v.DstToken,
		LpReceiverAddress: v.LpReceiverAddress,
		MsmqName:          v.MsmqName,
		SrcChainRawID:     v.SrcChainRawID,
		SrcToken:          v.SrcToken,
		WalletName:        v.WalletName,
		WalletID:          v.WalletID,
		EnableHedge:       v.EnableHedge,
	}
	if v.SrcTokenBalance != nil {
		res.SrcTokenBalance = *v.SrcTokenBalance
	}
	if v.DstTokenBalance != nil {
		res.DstTokenBalance = *v.DstTokenBalance
	}
	if v.SrcTokenDecimals != nil {
		res.SrcTokenDecimals = *v.SrcTokenDecimals
	}
	if v.DstTokenDecimals != nil {
		res.DstTokenDecimals = *v.DstTokenDecimals
	}
	if v.SrcTokenBalance == nil {
		res.SrcTokenBalance = "0"
	}
	if v.DstTokenBalance == nil {
		res.DstTokenBalance = "0"
	}
	if v.SrcTokenDecimals == nil {
		res.SrcTokenDecimals = 18
	}
	if v.DstTokenDecimals == nil {
		res.DstTokenDecimals = 18
	}

	return res
}

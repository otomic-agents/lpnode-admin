// Code generated by goa v3.11.0, DO NOT EDIT.
//
// bridgeConfig HTTP server encoders and decoders
//
// Command:
// $ goa gen admin-panel/design

package server

import (
	bridgeconfig "admin-panel/gen/bridge_config"
	"context"
	"io"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeBridgeCreateResponse returns an encoder for responses returned by the
// bridgeConfig bridgeCreate endpoint.
func EncodeBridgeCreateResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*bridgeconfig.BridgeCreateResult)
		enc := encoder(ctx, w)
		body := NewBridgeCreateResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeBridgeCreateRequest returns a decoder for requests sent to the
// bridgeConfig bridgeCreate endpoint.
func DecodeBridgeCreateRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body BridgeCreateRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateBridgeCreateRequestBody(&body)
		if err != nil {
			return nil, err
		}
		payload := NewBridgeCreateBridgeItem(&body)

		return payload, nil
	}
}

// EncodeBridgeListResponse returns an encoder for responses returned by the
// bridgeConfig bridgeList endpoint.
func EncodeBridgeListResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*bridgeconfig.BridgeListResult)
		enc := encoder(ctx, w)
		body := NewBridgeListResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// EncodeBridgeDeleteResponse returns an encoder for responses returned by the
// bridgeConfig bridgeDelete endpoint.
func EncodeBridgeDeleteResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*bridgeconfig.BridgeDeleteResult)
		enc := encoder(ctx, w)
		body := NewBridgeDeleteResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeBridgeDeleteRequest returns a decoder for requests sent to the
// bridgeConfig bridgeDelete endpoint.
func DecodeBridgeDeleteRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body BridgeDeleteRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateBridgeDeleteRequestBody(&body)
		if err != nil {
			return nil, err
		}
		payload := NewBridgeDeleteDeleteBridgeFilter(&body)

		return payload, nil
	}
}

// EncodeBridgeTestResponse returns an encoder for responses returned by the
// bridgeConfig bridgeTest endpoint.
func EncodeBridgeTestResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*bridgeconfig.BridgeTestResult)
		enc := encoder(ctx, w)
		body := NewBridgeTestResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeBridgeTestRequest returns a decoder for requests sent to the
// bridgeConfig bridgeTest endpoint.
func DecodeBridgeTestRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body BridgeTestRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		payload := NewBridgeTestPayload(&body)

		return payload, nil
	}
}

// marshalBridgeconfigListBridgeItemToListBridgeItemResponseBody builds a value
// of type *ListBridgeItemResponseBody from a value of type
// *bridgeconfig.ListBridgeItem.
func marshalBridgeconfigListBridgeItemToListBridgeItemResponseBody(v *bridgeconfig.ListBridgeItem) *ListBridgeItemResponseBody {
	if v == nil {
		return nil
	}
	res := &ListBridgeItemResponseBody{
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

	return res
}

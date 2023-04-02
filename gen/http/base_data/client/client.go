// Code generated by goa v3.11.0, DO NOT EDIT.
//
// baseData client HTTP transport
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

// Client lists the baseData service endpoint HTTP clients.
type Client struct {
	// ChainDataList Doer is the HTTP client used to make requests to the
	// chainDataList endpoint.
	ChainDataListDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the baseData service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		ChainDataListDoer:   doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// ChainDataList returns an endpoint that makes HTTP requests to the baseData
// service chainDataList server.
func (c *Client) ChainDataList() goa.Endpoint {
	var (
		decodeResponse = DecodeChainDataListResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildChainDataListRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ChainDataListDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("baseData", "chainDataList", err)
		}
		return decodeResponse(resp)
	}
}

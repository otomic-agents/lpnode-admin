// Code generated by goa v3.11.0, DO NOT EDIT.
//
// installCtrlPanel client HTTP transport
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

// Client lists the installCtrlPanel service endpoint HTTP clients.
type Client struct {
	// ListInstall Doer is the HTTP client used to make requests to the listInstall
	// endpoint.
	ListInstallDoer goahttp.Doer

	// InstallLpClient Doer is the HTTP client used to make requests to the
	// installLpClient endpoint.
	InstallLpClientDoer goahttp.Doer

	// UninstallLpClient Doer is the HTTP client used to make requests to the
	// uninstallLpClient endpoint.
	UninstallLpClientDoer goahttp.Doer

	// InstallDeployment Doer is the HTTP client used to make requests to the
	// installDeployment endpoint.
	InstallDeploymentDoer goahttp.Doer

	// UninstallDeployment Doer is the HTTP client used to make requests to the
	// uninstallDeployment endpoint.
	UninstallDeploymentDoer goahttp.Doer

	// UpdateDeployment Doer is the HTTP client used to make requests to the
	// updateDeployment endpoint.
	UpdateDeploymentDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the installCtrlPanel service
// servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		ListInstallDoer:         doer,
		InstallLpClientDoer:     doer,
		UninstallLpClientDoer:   doer,
		InstallDeploymentDoer:   doer,
		UninstallDeploymentDoer: doer,
		UpdateDeploymentDoer:    doer,
		RestoreResponseBody:     restoreBody,
		scheme:                  scheme,
		host:                    host,
		decoder:                 dec,
		encoder:                 enc,
	}
}

// ListInstall returns an endpoint that makes HTTP requests to the
// installCtrlPanel service listInstall server.
func (c *Client) ListInstall() goa.Endpoint {
	var (
		encodeRequest  = EncodeListInstallRequest(c.encoder)
		decodeResponse = DecodeListInstallResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildListInstallRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ListInstallDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("installCtrlPanel", "listInstall", err)
		}
		return decodeResponse(resp)
	}
}

// InstallLpClient returns an endpoint that makes HTTP requests to the
// installCtrlPanel service installLpClient server.
func (c *Client) InstallLpClient() goa.Endpoint {
	var (
		encodeRequest  = EncodeInstallLpClientRequest(c.encoder)
		decodeResponse = DecodeInstallLpClientResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildInstallLpClientRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.InstallLpClientDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("installCtrlPanel", "installLpClient", err)
		}
		return decodeResponse(resp)
	}
}

// UninstallLpClient returns an endpoint that makes HTTP requests to the
// installCtrlPanel service uninstallLpClient server.
func (c *Client) UninstallLpClient() goa.Endpoint {
	var (
		encodeRequest  = EncodeUninstallLpClientRequest(c.encoder)
		decodeResponse = DecodeUninstallLpClientResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUninstallLpClientRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UninstallLpClientDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("installCtrlPanel", "uninstallLpClient", err)
		}
		return decodeResponse(resp)
	}
}

// InstallDeployment returns an endpoint that makes HTTP requests to the
// installCtrlPanel service installDeployment server.
func (c *Client) InstallDeployment() goa.Endpoint {
	var (
		encodeRequest  = EncodeInstallDeploymentRequest(c.encoder)
		decodeResponse = DecodeInstallDeploymentResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildInstallDeploymentRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.InstallDeploymentDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("installCtrlPanel", "installDeployment", err)
		}
		return decodeResponse(resp)
	}
}

// UninstallDeployment returns an endpoint that makes HTTP requests to the
// installCtrlPanel service uninstallDeployment server.
func (c *Client) UninstallDeployment() goa.Endpoint {
	var (
		encodeRequest  = EncodeUninstallDeploymentRequest(c.encoder)
		decodeResponse = DecodeUninstallDeploymentResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUninstallDeploymentRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UninstallDeploymentDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("installCtrlPanel", "uninstallDeployment", err)
		}
		return decodeResponse(resp)
	}
}

// UpdateDeployment returns an endpoint that makes HTTP requests to the
// installCtrlPanel service updateDeployment server.
func (c *Client) UpdateDeployment() goa.Endpoint {
	var (
		encodeRequest  = EncodeUpdateDeploymentRequest(c.encoder)
		decodeResponse = DecodeUpdateDeploymentResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUpdateDeploymentRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UpdateDeploymentDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("installCtrlPanel", "updateDeployment", err)
		}
		return decodeResponse(resp)
	}
}

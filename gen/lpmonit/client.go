// Code generated by goa v3.11.0, DO NOT EDIT.
//
// lpmonit client
//
// Command:
// $ goa gen admin-panel/design

package lpmonit

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "lpmonit" service client.
type Client struct {
	AddScriptEndpoint    goa.Endpoint
	ListScriptEndpoint   goa.Endpoint
	DeleteScriptEndpoint goa.Endpoint
	RunScriptEndpoint    goa.Endpoint
	RunResultEndpoint    goa.Endpoint
}

// NewClient initializes a "lpmonit" service client given the endpoints.
func NewClient(addScript, listScript, deleteScript, runScript, runResult goa.Endpoint) *Client {
	return &Client{
		AddScriptEndpoint:    addScript,
		ListScriptEndpoint:   listScript,
		DeleteScriptEndpoint: deleteScript,
		RunScriptEndpoint:    runScript,
		RunResultEndpoint:    runResult,
	}
}

// AddScript calls the "add_script" endpoint of the "lpmonit" service.
func (c *Client) AddScript(ctx context.Context, p *AddScriptPayload) (res *AddScriptResult, err error) {
	var ires interface{}
	ires, err = c.AddScriptEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AddScriptResult), nil
}

// ListScript calls the "list_script" endpoint of the "lpmonit" service.
func (c *Client) ListScript(ctx context.Context) (res *ListScriptResult, err error) {
	var ires interface{}
	ires, err = c.ListScriptEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*ListScriptResult), nil
}

// DeleteScript calls the "delete_script" endpoint of the "lpmonit" service.
func (c *Client) DeleteScript(ctx context.Context, p *DeleteScriptPayload) (res *DeleteScriptResult, err error) {
	var ires interface{}
	ires, err = c.DeleteScriptEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DeleteScriptResult), nil
}

// RunScript calls the "run_script" endpoint of the "lpmonit" service.
func (c *Client) RunScript(ctx context.Context, p *RunScriptPayload) (res *RunScriptResult, err error) {
	var ires interface{}
	ires, err = c.RunScriptEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RunScriptResult), nil
}

// RunResult calls the "run_result" endpoint of the "lpmonit" service.
func (c *Client) RunResult(ctx context.Context, p *RunResultPayload) (res *RunResultResult, err error) {
	var ires interface{}
	ires, err = c.RunResultEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RunResultResult), nil
}

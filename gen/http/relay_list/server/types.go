// Code generated by goa v3.11.0, DO NOT EDIT.
//
// relayList HTTP server types
//
// Command:
// $ goa gen admin-panel/design

package server

import (
	relaylist "admin-panel/gen/relay_list"
)

// ListRelayResponseBody is the type of the "relayList" service "listRelay"
// endpoint HTTP response body.
type ListRelayResponseBody struct {
	Code    int64                             `form:"code" json:"code" xml:"code"`
	Result  []*RelayListRelayItemResponseBody `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string                           `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// RelayListRelayItemResponseBody is used to define fields on response body
// types.
type RelayListRelayItemResponseBody struct {
	ID           *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	Name         *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Profile      *string `form:"profile,omitempty" json:"profile,omitempty" xml:"profile,omitempty"`
	LpIDFake     *string `form:"lpIdFake,omitempty" json:"lpIdFake,omitempty" xml:"lpIdFake,omitempty"`
	LpNodeAPIKey *string `form:"lpNodeApiKey,omitempty" json:"lpNodeApiKey,omitempty" xml:"lpNodeApiKey,omitempty"`
	RelayAPIKey  *string `form:"relayApiKey,omitempty" json:"relayApiKey,omitempty" xml:"relayApiKey,omitempty"`
	RelayURI     *string `form:"relayUri,omitempty" json:"relayUri,omitempty" xml:"relayUri,omitempty"`
}

// NewListRelayResponseBody builds the HTTP response body from the result of
// the "listRelay" endpoint of the "relayList" service.
func NewListRelayResponseBody(res *relaylist.ListRelayResult) *ListRelayResponseBody {
	body := &ListRelayResponseBody{
		Code:    res.Code,
		Message: res.Message,
	}
	if res.Result != nil {
		body.Result = make([]*RelayListRelayItemResponseBody, len(res.Result))
		for i, val := range res.Result {
			body.Result[i] = marshalRelaylistRelayListRelayItemToRelayListRelayItemResponseBody(val)
		}
	}
	return body
}

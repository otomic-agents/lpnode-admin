// Code generated by goa v3.11.0, DO NOT EDIT.
//
// settings HTTP server types
//
// Command:
// $ goa gen admin-panel/design

package server

import (
	settings "admin-panel/gen/settings"

	goa "goa.design/goa/v3/pkg"
)

// SettingsRequestBody is the type of the "settings" service "settings"
// endpoint HTTP request body.
type SettingsRequestBody struct {
	RelayURI *string `form:"relayUri,omitempty" json:"relayUri,omitempty" xml:"relayUri,omitempty"`
}

// SettingsResponseBody is the type of the "settings" service "settings"
// endpoint HTTP response body.
type SettingsResponseBody struct {
	Code *int64 `form:"code,omitempty" json:"code,omitempty" xml:"code,omitempty"`
	// result
	Result  *int64  `form:"result,omitempty" json:"result,omitempty" xml:"result,omitempty"`
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// NewSettingsResponseBody builds the HTTP response body from the result of the
// "settings" endpoint of the "settings" service.
func NewSettingsResponseBody(res *settings.SettingsResult) *SettingsResponseBody {
	body := &SettingsResponseBody{
		Code:    res.Code,
		Result:  res.Result,
		Message: res.Message,
	}
	return body
}

// NewSettingsPayload builds a settings service settings endpoint payload.
func NewSettingsPayload(body *SettingsRequestBody) *settings.SettingsPayload {
	v := &settings.SettingsPayload{
		RelayURI: *body.RelayURI,
	}

	return v
}

// ValidateSettingsRequestBody runs the validations defined on
// SettingsRequestBody
func ValidateSettingsRequestBody(body *SettingsRequestBody) (err error) {
	if body.RelayURI == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("relayUri", "body"))
	}
	return
}

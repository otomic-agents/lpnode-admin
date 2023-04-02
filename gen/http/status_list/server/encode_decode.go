// Code generated by goa v3.11.0, DO NOT EDIT.
//
// statusList HTTP server encoders and decoders
//
// Command:
// $ goa gen admin-panel/design

package server

import (
	statuslist "admin-panel/gen/status_list"
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
)

// EncodeStatListResponse returns an encoder for responses returned by the
// statusList statList endpoint.
func EncodeStatListResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*statuslist.StatListResult)
		enc := encoder(ctx, w)
		body := NewStatListResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// marshalStatuslistStatusListItemToStatusListItemResponseBody builds a value
// of type *StatusListItemResponseBody from a value of type
// *statuslist.StatusListItem.
func marshalStatuslistStatusListItemToStatusListItemResponseBody(v *statuslist.StatusListItem) *StatusListItemResponseBody {
	if v == nil {
		return nil
	}
	res := &StatusListItemResponseBody{
		InstallType: v.InstallType,
		StatusKey:   v.StatusKey,
		StatusBody:  v.StatusBody,
		Name:        v.Name,
		ErrMessage:  v.ErrMessage,
	}

	return res
}

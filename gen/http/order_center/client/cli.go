// Code generated by goa v3.11.0, DO NOT EDIT.
//
// orderCenter HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	ordercenter "admin-panel/gen/order_center"
	"encoding/json"
	"fmt"
)

// BuildListPayload builds the payload for the orderCenter list endpoint from
// CLI flags.
func BuildListPayload(orderCenterListBody string) (*ordercenter.ListPayload, error) {
	var err error
	var body ListRequestBody
	{
		err = json.Unmarshal([]byte(orderCenterListBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"page\": 296079324223575006,\n      \"pageSize\": 5768036956203467588,\n      \"status\": 7455122834547886855\n   }'")
		}
	}
	v := &ordercenter.ListPayload{
		Status:   body.Status,
		Page:     body.Page,
		PageSize: body.PageSize,
	}

	return v, nil
}

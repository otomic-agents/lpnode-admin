// Code generated by goa v3.11.0, DO NOT EDIT.
//
// bridgeConfig HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	bridgeconfig "admin-panel/gen/bridge_config"
	"encoding/json"
	"fmt"
)

// BuildBridgeCreatePayload builds the payload for the bridgeConfig
// bridgeCreate endpoint from CLI flags.
func BuildBridgeCreatePayload(bridgeConfigBridgeCreateBody string) (*bridgeconfig.BridgeItem, error) {
	var err error
	var body BridgeCreateRequestBody
	{
		err = json.Unmarshal([]byte(bridgeConfigBridgeCreateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"ammName\": \"Possimus fuga in.\",\n      \"bridgeName\": \"Ut ipsam.\",\n      \"dstChainId\": \"Doloremque voluptatem id eligendi aut reprehenderit.\",\n      \"dstTokenId\": \"Odio totam.\",\n      \"enableHedge\": false,\n      \"enableLimiter\": true,\n      \"srcChainId\": \"Sequi laudantium dolore non ullam et.\",\n      \"srcTokenId\": \"Reiciendis odio voluptas quas alias dolorum quae.\",\n      \"srcWalletId\": \"Quod sint.\",\n      \"walletId\": \"Nesciunt exercitationem voluptatem sint.\"\n   }'")
		}
	}
	v := &bridgeconfig.BridgeItem{
		BridgeName:    body.BridgeName,
		SrcChainID:    body.SrcChainID,
		DstChainID:    body.DstChainID,
		SrcTokenID:    body.SrcTokenID,
		DstTokenID:    body.DstTokenID,
		WalletID:      body.WalletID,
		SrcWalletID:   body.SrcWalletID,
		AmmName:       body.AmmName,
		EnableHedge:   body.EnableHedge,
		EnableLimiter: body.EnableLimiter,
	}
	{
		var zero bool
		if v.EnableHedge == zero {
			v.EnableHedge = true
		}
	}
	{
		var zero bool
		if v.EnableLimiter == zero {
			v.EnableLimiter = true
		}
	}

	return v, nil
}

// BuildBridgeDeletePayload builds the payload for the bridgeConfig
// bridgeDelete endpoint from CLI flags.
func BuildBridgeDeletePayload(bridgeConfigBridgeDeleteBody string) (*bridgeconfig.DeleteBridgeFilter, error) {
	var err error
	var body BridgeDeleteRequestBody
	{
		err = json.Unmarshal([]byte(bridgeConfigBridgeDeleteBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"id\": \"Aliquam sed hic ut.\"\n   }'")
		}
	}
	v := &bridgeconfig.DeleteBridgeFilter{
		ID: body.ID,
	}

	return v, nil
}

// BuildBridgeTestPayload builds the payload for the bridgeConfig bridgeTest
// endpoint from CLI flags.
func BuildBridgeTestPayload(bridgeConfigBridgeTestBody string) (*bridgeconfig.BridgeTestPayload, error) {
	var err error
	var body BridgeTestRequestBody
	{
		err = json.Unmarshal([]byte(bridgeConfigBridgeTestBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"id\": \"Ea et ut qui voluptatem expedita.\"\n   }'")
		}
	}
	v := &bridgeconfig.BridgeTestPayload{
		ID: body.ID,
	}

	return v, nil
}

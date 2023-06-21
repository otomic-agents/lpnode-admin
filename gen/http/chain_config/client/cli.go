// Code generated by goa v3.11.0, DO NOT EDIT.
//
// chainConfig HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	chainconfig "admin-panel/gen/chain_config"
	"encoding/json"
	"fmt"
)

// BuildSetChainListPayload builds the payload for the chainConfig setChainList
// endpoint from CLI flags.
func BuildSetChainListPayload(chainConfigSetChainListBody string) (*chainconfig.SetChainListPayload, error) {
	var err error
	var body SetChainListRequestBody
	{
		err = json.Unmarshal([]byte(chainConfigSetChainListBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"chainList\": [\n         {\n            \"chainId\": 8039505103587176397,\n            \"chainName\": \"Repellendus eius deserunt quaerat repudiandae.\",\n            \"name\": \"Sint exercitationem.\",\n            \"tokenName\": \"Praesentium quo ipsam atque cumque.\"\n         },\n         {\n            \"chainId\": 8039505103587176397,\n            \"chainName\": \"Repellendus eius deserunt quaerat repudiandae.\",\n            \"name\": \"Sint exercitationem.\",\n            \"tokenName\": \"Praesentium quo ipsam atque cumque.\"\n         }\n      ]\n   }'")
		}
	}
	v := &chainconfig.SetChainListPayload{}
	if body.ChainList != nil {
		v.ChainList = make([]*chainconfig.ChainDataItem, len(body.ChainList))
		for i, val := range body.ChainList {
			v.ChainList[i] = marshalChainDataItemRequestBodyToChainconfigChainDataItem(val)
		}
	}

	return v, nil
}

// BuildDelChainListPayload builds the payload for the chainConfig delChainList
// endpoint from CLI flags.
func BuildDelChainListPayload(chainConfigDelChainListBody string) (*chainconfig.DelChainListPayload, error) {
	var err error
	var body DelChainListRequestBody
	{
		err = json.Unmarshal([]byte(chainConfigDelChainListBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"_id\": \"Impedit nostrum eos et.\",\n      \"chainId\": 2031005260904629279\n   }'")
		}
	}
	v := &chainconfig.DelChainListPayload{
		ChainID: body.ChainID,
		ID:      body.ID,
	}

	return v, nil
}

// BuildSetChainGasUsdPayload builds the payload for the chainConfig
// setChainGasUsd endpoint from CLI flags.
func BuildSetChainGasUsdPayload(chainConfigSetChainGasUsdBody string) (*chainconfig.SetChainGasUsdPayload, error) {
	var err error
	var body SetChainGasUsdRequestBody
	{
		err = json.Unmarshal([]byte(chainConfigSetChainGasUsdBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"_id\": \"Doloremque dolores est quas reprehenderit vel esse.\",\n      \"chainId\": 2885005782751723716,\n      \"usd\": 5450437637379237693\n   }'")
		}
	}
	v := &chainconfig.SetChainGasUsdPayload{
		ChainID: body.ChainID,
		ID:      body.ID,
		Usd:     body.Usd,
	}

	return v, nil
}

// BuildSetChainClientConfigPayload builds the payload for the chainConfig
// setChainClientConfig endpoint from CLI flags.
func BuildSetChainClientConfigPayload(chainConfigSetChainClientConfigBody string) (*chainconfig.SetChainClientConfigPayload, error) {
	var err error
	var body SetChainClientConfigRequestBody
	{
		err = json.Unmarshal([]byte(chainConfigSetChainClientConfigBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"chainData\": \"Et quidem pariatur sunt et corrupti.\",\n      \"chainId\": 1837657479056530763\n   }'")
		}
	}
	v := &chainconfig.SetChainClientConfigPayload{
		ChainID:   body.ChainID,
		ChainData: body.ChainData,
	}

	return v, nil
}

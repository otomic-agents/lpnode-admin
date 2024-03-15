package design

import (
	. "goa.design/goa/v3/dsl"
)

var chainConfig_ChainDataItem = Type("chainDataItem", func() {
	Attribute("chainId", Int64, "ChainId")
	Attribute("chainName", String, "chainName")
	Attribute("name", String, "")
	Attribute("tokenName", String, "tokenName")
})
var chainConfig_walletDataItem = Type("chainWallet", func() {
	Attribute("walletName", String, "")
	Attribute("privateKey", String, "")
	Attribute("tokenList", ArrayOf(String), "token list")
})
var _ = Service("chainConfig", func() {
	Description("used to configure basic chain settings")
	Method("setChainList", func() {
		Payload(func() {
			Attribute("chainList", ArrayOf(chainConfig_ChainDataItem), "listData")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(chainConfig_ChainDataItem), "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/chainConfig/setChainList")
		})
	})
	Method("delChainList", func() {
		Description("used to delete basic data for a chain")
		Payload(func() {
			Attribute("chainId", Int64, "ChainId")
			Attribute("_id", String, "mongodb id")
			Required("_id", "chainId")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", Int64, "")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/chainConfig/delChainList")
		})
	})
	Method("chainList", func() {
		Description("list chain and append chain service status, like client runtime status")
		Payload(func() {
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(chainConfig_ChainDataItem), "chain list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/chainConfig/chainList")
		})
	})
	Method("setChainGasUsd", func() {
		Payload(func() {
			Attribute("chainId", Int64, "ChainId")
			Attribute("_id", String, "mongodb id")
			Attribute("usd", Int64, "usd value")
			Required("usd", "_id", "chainId")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", Int64, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/chainConfig/setChainGasUsd")
		})
	})
	Method("setChainClientConfig", func() {
		Payload(func() {
			Attribute("chainId", Int64, "ChainId")
			Attribute("chainData", String, "JSON Stored String")
			Required("chainId", "chainData")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", Int64, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/chainConfig/setChainClientConfig")
		})
	})
})

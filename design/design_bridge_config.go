package design

import (
	. "goa.design/goa/v3/dsl"
)

var BridgeConfig_bridgeItem = Type("bridgeItem", func() {
	// Attribute("installId", String, "id after service installed")
	Attribute("bridgeName", String, "bridge name ****")
	Attribute("srcChainId", String, "mongodb primary key, from basedata")
	Attribute("dstChainId", String, "mongodb primary key, from basedata")
	Attribute("srcTokenId", String, "mongodb primary key, from tokenlist")
	Attribute("dstTokenId", String, "mongodb primary key, from tokenlist")
	Attribute("walletId", String, "mongodb primary key, from walletlist")    // payment wallet info
	Attribute("srcWalletId", String, "mongodb primary key, from walletlist") // receiving wallet info
	Attribute("ammName", String, "amm name at install")
	Attribute("relayApiKey", String, "relay api key")
	Attribute("relayUri", String, "relayUri")
	Attribute("enableHedge", Boolean, func() { // enable hedging for token pair
		Default(true)
	})
	Attribute("enableLimiter", Boolean, func() { // enable trade limit
		Default(true)
	})
	Required("bridgeName", "srcChainId", "dstChainId", "srcTokenId", "dstTokenId", "srcWalletId", "walletId", "ammName", "relayApiKey","relayUri")
})
var BridgeConfig_listItem = Type("listBridgeItem", func() {
	Attribute("_id", String)
	Attribute("dstChainId", String)
	Attribute("dstTokenId", String)
	Attribute("srcChainId", String)
	Attribute("srcTokenId", String)
	Attribute("ammName", String)
	Attribute("bridgeName", String)
	Attribute("dstChainRawId", Int64)
	Attribute("dstClientUri", String)
	Attribute("dstToken", String)
	Attribute("lpReceiverAddress", String)
	Attribute("msmqName", String)
	Attribute("srcChainRawId", Int64)
	Attribute("srcToken", String)
	Attribute("walletName", String)
	Attribute("walletId", String)
	Attribute("enableHedge", Boolean)
})
var BridgeConfig_DeleteBridgeFilter = Type("deleteBridgeFilter", func() {
	Attribute("id", String, "mongodb primary key")
	Required("id")
})
var _ = Service("bridgeConfig", func() {
	Method("bridgeCreate", func() {
		Description("used to create cross-chain config")
		Payload(BridgeConfig_bridgeItem)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/bridge/create")
		})
	})
	Method("bridgeList", func() {
		Payload(func() {
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(BridgeConfig_listItem), "chain list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/bridge/list")
		})
	})
	Method("bridgeDelete", func() {
		Payload(BridgeConfig_DeleteBridgeFilter)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/bridge/delete")
		})
	})
	Method("bridgeTest", func() {
		Payload(func() {
			Attribute("id", String)
		})
		Result(func() {
			Attribute("code", Int64, "")

		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/bridge/test")
		})
	})
})

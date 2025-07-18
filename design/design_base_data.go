package design

import (
	. "goa.design/goa/v3/dsl"
)

var baseData_ChainDataItem = Type("ChainDataItem", func() {
	Attribute("id", String, "Chain ID in the database")
	Attribute("chainId", Int64, "Chain ID")
	Attribute("name", String, "chain name")
	Attribute("chainName", String, "full chain name")
	Attribute("tokenName", String, "token name")
})
var baseData_Lpinfo = Type("LpInfo", func() {
	Attribute("name", String)
	Attribute("profile", String)
})
var baseData_WalletTokenItem = Type("WalletTokenItem", func() {
	Attribute("address", String, "token address")
	Attribute("symbol", String, "token symbol")
	Attribute("decimals", Int32, "token decimals")
})

var baseData_WalletItem = Type("WalletItem", func() {
	Attribute("wallet_name", String, "wallet name")
	Attribute("address", String, "wallet address")
	Attribute("can_sign", Boolean, "whether can sign")
	Attribute("can_sign_712", Boolean, "whether can sign 712")
	Attribute("type", String, "wallet type")
	Attribute("signature_service_address", String, "signature service address")
	Attribute("tokens", ArrayOf(baseData_WalletTokenItem), "token list")
})
var _ = Service("baseData", func() {
	Description("used to manage basic data")
	Method("chainDataList", func() {
		Description("used to return basic chain data")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(baseData_ChainDataItem), "list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/baseData/chainDataList")
		})
	})
	Method("getLpInfo", func() {
		Description("used to return basic chain data")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", baseData_Lpinfo, "")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/baseData/getLpInfo")
		})
	})
	Method("runTimeEnv", func() {
		Description("used to return runtime environment")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", String, "list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/baseData/runTimeEnv")
		})
	})
	Method("getWalletAndTokens", func() {
		Description("Get wallet list with their associated tokens")
		Payload(func() {
			Attribute("chainId", Int64, "Chain ID")
			Required("chainId")
		})
		Result(func() {
			Attribute("code", Int64, "response code")
			Attribute("result", ArrayOf(baseData_WalletItem), "wallet list with tokens")
			Attribute("message", String, "response message")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/baseData/wallets")
		})
	})
})

package design

import (
	. "goa.design/goa/v3/dsl"
)

var TokenManager_DeleteTokenFilter = Type("DeleteTokenFilter", func() {
	Attribute("_id", String, "mongodb primary key")
	Required("_id")
})
var TokenManager_TokenItem = Type("TokenItem", func() {
	Attribute("_id", String)
	Attribute("tokenId", String) // near only
	Attribute("chainId", Int64, "chain id")
	Attribute("address", String, "token address") // address is receipt_id on near
	Attribute("tokenName", String, "name in token contract")
	Attribute("marketName", String, "corresponding name in cex")
	Attribute("precision", Int64, func() {
		Maximum(18)
		Minimum(6)
	})
	Attribute("chainType", String, "")
	Attribute("coinType", String, func() {
		Enum("stable_coin", "coin")
	})
	Required("chainId", "address", "marketName", "precision", "coinType")
})

var _ = Service("tokenManager", func() {
	Description("used to manage all tokens")
	Method("tokenList", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(TokenManager_TokenItem), "list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/tokenManager/list")
		})
	})
	Method("tokenCreate", func() {
		Payload(TokenManager_TokenItem)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/tokenManager/create")
		})
	})
	Method("tokenDelete", func() {
		Payload(TokenManager_DeleteTokenFilter)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/tokenManager/delete")
		})
	})
})

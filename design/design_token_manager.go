package design

import (
	. "goa.design/goa/v3/dsl"
)

var TokenManager_DeleteTokenFilter = Type("DeleteTokenFilter", func() {
	Attribute("_id", String, "Mongodb 的主键")
	Required("_id")
})
var TokenManager_TokenItem = Type("TokenItem", func() {
	Attribute("_id", String)
	Attribute("tokenId", String) // near 多一个这个
	Attribute("chainId", Int64, "链Id")
	Attribute("address", String, "token address") // near 时候 address 为 receipt_id
	Attribute("tokenName", String, "token 合约中的名称")
	Attribute("marketName", String, "Cex中所对应的名称")
	Attribute("precision", Int64, func() {
		Maximum(18)
		Minimum(6)
	})
	Attribute("chainType", String, "") // 后端根据chanID自动填充
	Attribute("coinType", String, func() {
		Enum("stable_coin", "coin")
	})
	Required("chainId", "address", "marketName", "precision", "coinType")
})

var _ = Service("tokenManager", func() {
	Description("用于管理所有的token")
	Method("tokenList", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(TokenManager_TokenItem), "添加成功的链")
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
			Attribute("result", Int64, "创建影响的行数")
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
			Attribute("result", Int64, "是否删除成功")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/tokenManager/delete")
		})
	})
})

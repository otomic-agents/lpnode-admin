package design

import (
	. "goa.design/goa/v3/dsl"
)

var chainConfig_ChainDataItem = Type("chainDataItem", func() {
	Attribute("chainId", Int64, "ChainId")
	Attribute("chainName", String, "chainName") // 链的全称
	Attribute("name", String, "")               //链的简称
	Attribute("tokenName", String, "tokenName") //链token币的名称， 对冲，或者获取行情的时候使用
})
var chainConfig_walletDataItem = Type("chainWallet", func() {
	Attribute("walletName", String, "")
	Attribute("privateKey", String, "")
	Attribute("tokenList", ArrayOf(String), "Token列表")
})
var _ = Service("chainConfig", func() {
	Description("用于配置chain的基础设置")
	Method("setChainList", func() {
		Description("用于配置chain的基础设置,批量设置接口Upsert")
		Payload(func() {
			Attribute("chainList", ArrayOf(chainConfig_ChainDataItem), "listData")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(chainConfig_ChainDataItem), "添加成功的链")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/chainConfig/setChainList")
		})
	})
	Method("delChainList", func() {
		Description("用于删除一项链的基础设置")
		Payload(func() {
			Attribute("chainId", Int64, "ChainId")
			Attribute("_id", String, "mongodb的id")
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
		Description("列出链的列表，并附加链相关服务的状态，如Client 运行时状态")
		Payload(func() {
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(chainConfig_ChainDataItem), "链的列表")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/chainConfig/chainList")
		})
	})
	Method("setChainGasUsd", func() {
		Description("用于换目标链原生币种时，最少换多少USD价值的原生币")
		Payload(func() {
			Attribute("chainId", Int64, "ChainId")
			Attribute("_id", String, "mongodb的id")
			Attribute("usd", Int64, "usd value")
			Required("usd", "_id", "chainId")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", Int64, "是否成功")
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
			Attribute("data", Int64, "是否成功")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/chainConfig/setChainClientConfig")
		})
	})
})

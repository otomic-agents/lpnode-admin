package design

import (
	. "goa.design/goa/v3/dsl"
)

var BridgeConfig_bridgeItem = Type("bridgeItem", func() {
	// Attribute("installId", String, "安装服务后的安装id")
	Attribute("bridgeName", String, "bridge的Name ****")
	Attribute("srcChainId", String, "mongodb的主键,baseData中获取")
	Attribute("dstChainId", String, "mongodb的主键,baseData中获取")
	Attribute("srcTokenId", String, "mongodb的主键,tokenList中获取")
	Attribute("dstTokenId", String, "mongodb的主键,tokenList中获取")
	Attribute("walletId", String, "mongodb的主键,walletList 中获取")    // 付款钱包信息
	Attribute("srcWalletId", String, "mongodb的主键,walletList 中获取") // 收款钱包信息
	Attribute("ammName", String, "amm安装时候的name")
	Attribute("enableHedge", Boolean, func() { // 是否开启币对的对冲
		Default(true)
	})
	//Attribute("receiveAddress", String, "") // 原始链收款地址 , 这个应该程序，通过收款钱包 id来获取
	// Attribute("msmqName", String, "应当根据链和TokenAddress自动生成一个 ****")
	// Attribute("dstChainClientUri", String, "根据dstChainId 查install记录 并找到对应的serviceName ****")
	// Attribute("relay_api_key", String, "***** ") // 最后再处理这个逻辑
	Required("bridgeName", "srcChainId", "dstChainId", "srcTokenId", "dstTokenId", "srcWalletId", "walletId", "ammName")
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
	Attribute("id", String, "Mongodb 的主键")
	Required("id")
})
var _ = Service("bridgeConfig", func() {
	Method("bridgeCreate", func() {
		Description("用于创建跨链配置")
		Payload(BridgeConfig_bridgeItem)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "是否成功")
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
			Attribute("result", ArrayOf(BridgeConfig_listItem), "链的列表")
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
			Attribute("result", Int64, "是否删除成功")
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

package design

import (
	. "goa.design/goa/v3/dsl"
)

var DexWallet_WalletRow = Type("walletRow", func() {
	Attribute("id", String, "mongodb主键")
	Attribute("walletName", String, "")
	Attribute("privateKey", String, "")
	Attribute("address", String, "")
	Attribute("chainType", String, "")                  // evm  near
	Attribute("accountId", String, "wallet对应的人类可阅读的名称") // Near 有这个玩意
	Attribute("chainId", Int64, "链的Id")
	Attribute("storeId", String, "")
	Attribute("vaultHostType", String)
	Attribute("vaultName", String)
	Attribute("vaultSecertType", String)
	Attribute("walletType", String, func() {
		Enum("privateKey", "storeId")
	})
	Required("walletName", "chainId", "chainType", "walletType")
})
var DexWallet_VaultRow = Type("vaultRow", func() {
	Attribute("address", String, "地址")
	Attribute("hostType", String, "托管类型")
	Attribute("id", String, "存储Id")
	Attribute("name", String, "钱包名称")
	Attribute("secertType", String, "私钥类型")
})
var DexWallet_DeleteFilter = Type("deleteFilter", func() {
	Attribute("id", String, "Mongodb 的主键")
	Required("id")
})
var _ = Service("dexWallet", func() {
	Description("用于管理账号")
	Method("listDexWallet", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(DexWallet_WalletRow), "钱包别表")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/dexWallet/list")
		})
	})
	Method("createDexWallet", func() {
		Payload(DexWallet_WalletRow)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", func() {
				Attribute("_id", String) // Mongodb的id
			})
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/dexWallet/create")
		})
	})
	Method("deleteDexWallet", func() {
		Payload(DexWallet_DeleteFilter)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "是否删除成功")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/dexWallet/delete")
		})
	})
	Method("vaultList", func() {
		Payload(func() {})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(DexWallet_VaultRow), "列表")
			Attribute("message", String)
			Required("result")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/dexWallet/vaultList")
		})
	})
})

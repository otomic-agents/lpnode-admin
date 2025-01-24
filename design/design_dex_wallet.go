package design

import (
	. "goa.design/goa/v3/dsl"
)

var DexWallet_WalletRow = Type("walletRow", func() {
	Attribute("id", String, "mongodb primary key")
	Attribute("walletName", String, "")
	Attribute("privateKey", String, "")
	Attribute("address", String, "")
	Attribute("chainType", String, "")       // evm  near
	Attribute("accountId", String, "wallet") // Near only
	Attribute("chainId", Int64, "chain Id")
	Attribute("storeId", String, "")
	Attribute("vaultHostType", String)
	Attribute("vaultName", String)
	Attribute("vaultSecertType", String)
	Attribute("signServiceEndpoint",String)
	Attribute("walletType", String, func() {
		Enum("privateKey", "storeId")
	})
	Attribute("balance",String)
	Required("walletName", "chainId", "chainType", "walletType")
})
var DexWallet_VaultRow = Type("vaultRow", func() {
	Attribute("address", String, "address")
	Attribute("hostType", String, "host type")
	Attribute("id", String, "storeId")
	Attribute("name", String, "store name")
	Attribute("secertType", String, "store secert type")
})
var DexWallet_DeleteFilter = Type("deleteFilter", func() {
	Attribute("id", String, "mongodb primary key")
	Required("id")
})
var _ = Service("dexWallet", func() {
	Description("used to manage wallets")
	Method("listDexWallet", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(DexWallet_WalletRow), "wallet list")
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
				Attribute("_id", String) // Mongodb id
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
			Attribute("result", Int64, "result")
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
			Attribute("result", ArrayOf(DexWallet_VaultRow), "list")
			Attribute("message", String)
			Required("result")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/dexWallet/vaultList")
		})
	})
	Method("updateLpWallet", func() {
		Payload(func() {})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", String, "list")
			Attribute("message", String)
			Required("result")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/dexWallet/updateLpWallet")
		})

	})
})

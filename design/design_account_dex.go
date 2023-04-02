package design

import (
	. "goa.design/goa/v3/dsl"
)

var accountDex_DexAccountBalance = Type("DexAccountBalance", func() {
	Attribute("token", String, "")
	Attribute("tokenName", String, "")
	Attribute("amount", String, "")
	Attribute("free", String, "")
	Attribute("locked", String, "")
	Attribute("price", String)
})
var _ = Service("accountDex", func() {
	Method("walletInfo", func() {
		Payload(func() {
			Attribute("chainId", Int64, "链的Id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(accountDex_DexAccountBalance), "是否成功")
			Attribute("message", String)
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/account/dex/walletInfo")
		})
	})
})

package design

import (
	. "goa.design/goa/v3/dsl"
)

var accountCex_CexAccountBalance = Type("CexAccountBalance", func() {
	Attribute("asset", String, "")
	Attribute("total", String, "")
	Attribute("free", String, "")
	Attribute("locked", String, "")
	Attribute("price", String)
})
var _ = Service("accountCex", func() {
	Method("walletInfo", func() {
		Payload(func() {
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(accountCex_CexAccountBalance), "是否成功")
			Attribute("message", String)
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/account/cex/walletInfo")
		})
	})
})

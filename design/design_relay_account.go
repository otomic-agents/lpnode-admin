package design

import (
	. "goa.design/goa/v3/dsl"
)

var relayAccount_AccountItem = Type("relayAccountItem", func() {
	Attribute("id", String)
	Attribute("name", String)
	Attribute("profile", String)
	Attribute("lpIdFake", String)
	Attribute("lpNodeApiKey", String)
	Attribute("relayApiKey", String)
})
var _ = Service("relayAccount", func() {
	Description("用于管理Lp和relay之间的账号")
	Method("listAccount", func() {
		Payload(func() {})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", ArrayOf(relayAccount_AccountItem))
			Attribute("message", String)
			Required("code")
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/relayAccount/list")
		})
	})
	Method("registerAccount", func() {
		Payload(func() {
			Attribute("name", String)
			Attribute("profile", String)
			Required("name")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", relayAccount_AccountItem)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/relayAccount/register")
		})
	})
})

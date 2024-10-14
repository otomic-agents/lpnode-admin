package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("relayAccount", func() {
	Description("used to manage lp account on relay")
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
			// Attribute("name", String)
			Attribute("relayUrl",String)
			Attribute("profile", String)
			Required("relayUrl","profile")
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
	Method("deleteAccount", func() {
		Payload(func() {
			Attribute("id", String)
			Required("id")
		})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", String)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/relayAccount/delete")
		})
	})
})
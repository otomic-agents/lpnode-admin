package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("lpRegister", func() {
	Description("用于管理Lp到Client的注册")

	Method("registerAll", func() {
		Payload(func() {})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", String)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/relayAccount/registerAll")
		})
	})
	Method("unRegisterAll", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", String)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/lpRegister/unRegisterAll")
		})
	})
})

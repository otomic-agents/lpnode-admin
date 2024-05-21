package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("settings", func() {
	Description("used to manage order")
	Method("settings", func() {
		Payload(func() {
			Attribute("relayUri", String)
			Required("relayUri")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/settings/save")
		})
	})
})

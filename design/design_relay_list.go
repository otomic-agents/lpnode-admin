package design

import (
	. "goa.design/goa/v3/dsl"
)

var relayList_RelayItem = Type("relayListRelayItem", func() {
	Attribute("id", String)
	Attribute("name", String)
	Attribute("profile", String)
	Attribute("lpIdFake", String)
	Attribute("lpNodeApiKey", String)
	Attribute("relayApiKey", String)
	Attribute("relayUri", String)
})
var _ = Service("relayList", func() {
	Description("used to manage lp account on relay")
	Method("listRelay", func() {
		Payload(func() {})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", ArrayOf(relayList_RelayItem))
			Attribute("message", String)
			Required("code")
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/relayList/list")
		})
	})
})

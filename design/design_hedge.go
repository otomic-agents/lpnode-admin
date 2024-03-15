package design

import (
	. "goa.design/goa/v3/dsl"
)

var hedge_Hedge = Type("HedgeItem", func() {
	Attribute("id", String, "")
	Attribute("hedgeType", String, "")
})
var _ = Service("hedge", func() {
	Description("hedge basic configuration")
	Method("list", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(hedge_Hedge), "list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/hedge/list")
		})
	})
	Method("edit", func() {
		Payload(func() {
			Attribute("hedge", hedge_Hedge, "")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", Int64, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/hedge/edit")
		})
	})
	Method("del", func() {
		Payload(func() {
			Attribute("id", String, "primary key")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", Int64, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/hedge/del")
		})
	})
})

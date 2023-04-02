package design

import (
	. "goa.design/goa/v3/dsl"
)

var hedge_Hedge = Type("HedgeItem", func() {
	Attribute("id", String, "")
	Attribute("hedgeType", String, "")
})
var _ = Service("hedge", func() {
	Description("对冲的基本配置")
	Method("list", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(hedge_Hedge), "添加成功的链")
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
			Attribute("data", Int64, "添加成功的链")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/hedge/edit")
		})
	})
	Method("del", func() {
		Payload(func() {
			Attribute("id", String, "删除的主键Id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", Int64, "是否成功删除")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/hedge/del")
		})
	})
})

package design

import (
	. "goa.design/goa/v3/dsl"
)

var ammOrderCenter_OrderItem = Type("OrderCenterAmmOrderItem", Any)
var ammOrderCenterRetResult = Type("ammOrderCenterRetResult", func() {
	Attribute("list", ArrayOf(ammOrderCenter_OrderItem))
	Attribute("pageCount", Int64)
})
var _ = Service("ammOrderCenter", func() {
	Description("used to manage amm order")
	Method("list", func() {
		Payload(func() {
			Attribute("status", Int64)
			Attribute("ammName", String)
			Attribute("page", Int64)
			Attribute("pageSize", Int64)
			Required("ammName")
		})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", ammOrderCenterRetResult)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ammOrderCenter/list")
		})
	})
})

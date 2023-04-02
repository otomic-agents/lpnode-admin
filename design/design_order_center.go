package design

import (
	. "goa.design/goa/v3/dsl"
)

var orderCenter_OrderItem = Type("OrderCenterOrderItem", Any)
var orderCenterRetResult = Type("orderCenterRetResult", func() {
	Attribute("list", ArrayOf(orderCenter_OrderItem))
	Attribute("pageCount", Int64)
})
var _ = Service("orderCenter", func() {
	Description("用于管理orderCenter")
	Method("list", func() {
		Payload(func() {
			Attribute("status", Int64)
			Attribute("page", Int64)
			Attribute("pageSize", Int64)
		})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", orderCenterRetResult)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/orderCenter/list")
		})
	})
})

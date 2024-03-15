package design

import (
	. "goa.design/goa/v3/dsl"
)

var statusList_ListItem = Type("StatusListItem", func() {
	Attribute("installType", String)
	Attribute("statusKey", String)
	Attribute("statusBody", String)
	Attribute("name", String)
	Attribute("errMessage", String)
})
var _ = Service("statusList", func() {
	Description("used to manage install status")
	Method("statList", func() {
		Payload(func() {})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", ArrayOf(statusList_ListItem))
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/staustList/list")
		})
	})
})

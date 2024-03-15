package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("authenticationLimiter", func() {
	Description("used to manage ordercenter")

	Method("getAuthenticationLimiter", func() {
		Description("used to query limiter information")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", String, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/authenticationLimiter/getAuthenticationLimiter")
		})
	})
	Method("setAuthenticationLimiter", func() {
		Description("set limit information")
		Payload(func() {
			Attribute("authenticationLimiter", String)
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", String, "result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/authenticationLimiter/setAuthenticationLimiter")
		})
	})
	Method("delAuthenticationLimiter", func() {
		Description("delete system limit information")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/authenticationLimiter/delAuthenticationLimiter")
		})
	})
})

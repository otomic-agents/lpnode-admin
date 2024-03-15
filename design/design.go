package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("adminApiService", func() {
	Title("adminApiService")
	Description("Service for multiplying numbers, a Goa teaser")
	Server("apiService", func() {
		Host("0.0.0.0", func() {
			URI("http://0.0.0.0:18006")
			URI("grpc://0.0.0.0:18007")
		})
	})
})
var _ = Service("mainLogic", func() {
	Method("mainLogic", func() {
		Result(func() {
			// Attribute("code", Int64, "")
		})
		HTTP(func() {
			GET("/")
		})
	})
	Method("mainLogicLink", func() {
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("data", String, "result")
		})
		HTTP(func() {
			GET("/link")
		})
	})
})

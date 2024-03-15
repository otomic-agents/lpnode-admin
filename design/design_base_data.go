package design

import (
	. "goa.design/goa/v3/dsl"
)

var baseData_ChainDataItem = Type("ChainDataItem", func() {
	Attribute("id", String, "chain id in the database")
	Attribute("chainId", Int64, "chain Id")
	Attribute("name", String, "chain name")
	Attribute("chainName", String, "full chain name")
	Attribute("tokenName", String, "token name")
})
var _ = Service("baseData", func() {
	Description("used to manage basic data")
	Method("chainDataList", func() {
		Description("used to return basic chain data")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(baseData_ChainDataItem), "list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/baseData/chainDataList")
		})
	})
	Method("runTimeEnv", func() {
		Description("used to return runtime environment")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", String, "list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/baseData/runTimeEnv")
		})
	})
})

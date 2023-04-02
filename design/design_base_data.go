package design

import (
	. "goa.design/goa/v3/dsl"
)

var baseData_ChainDataItem = Type("ChainDataItem", func() {
	Attribute("id", String, "链在数据库中的id")
	Attribute("chainId", Int64, "链的Id")
	Attribute("name", String, "链名称")
	Attribute("chainName", String, "链全称")
	Attribute("tokenName", String, "Token币的名称")
})
var _ = Service("baseData", func() {
	Description("用于管理基础数据")
	Method("chainDataList", func() {
		Description("用于返回最基础的链的数据")
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(baseData_ChainDataItem), "列表")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/baseData/chainDataList")
		})
	})
})

package design

import (
	. "goa.design/goa/v3/dsl"
)

var configResource_ConfigResultIdItem = Type("ConfigResultIdItem", func() {
	Attribute("id", String, "mongodb的id")
	Attribute("clientId", String, "客户端提交的id")
	Required("id")
})
var configResource_ConfigResult = Type("ResourceConfigResult", func() {
	Attribute("id", String)
	Attribute("templateResult", String)
	Attribute("template", String)
	Attribute("clientId", String)
	Attribute("appName", String)
	Attribute("version", String)
	Attribute("versionHash", String)
	Required("clientId")

})
var _ = Service("configResource", func() {
	Method("createResource", func() {
		Payload(func() {
			Attribute("appName", String)
			Attribute("version", String)
			Attribute("clientId", String)
			Attribute("template", String)
			Required("appName", "clientId")
		})
		Result(func() {
			Attribute("code", Int64, "0是成功")
			Attribute("result", configResource_ConfigResultIdItem)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/configResource/create")
		})
	})
	Method("getResource", func() {
		Payload(func() {
			Attribute("clientId", String)
			Required("clientId")
		})
		Result(func() {
			Attribute("code", Int64, "0是成功")
			Attribute("result", configResource_ConfigResult)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/configResource/get")
		})
	})
	Method("listResource", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(configResource_ConfigResult))
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/configResource/list")
		})
	})
	Method("deleteResult", func() {
		Payload(func() {

		})
		Result(func() {
			Attribute("code", Int64, "")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/configResource/delete")
		})
	})
	Method("editResult", func() {
		Payload(func() {
			Attribute("templateResult", String)
			Attribute("template", String)
			Attribute("clientId", String)
			Attribute("appName", String)
			Attribute("version", String)
			Attribute("versionHash", String)
			Required("clientId", "templateResult")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("message", String, "")
			Attribute("result", String, "修改影响的Id")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/configResource/edit")
		})
	})
})

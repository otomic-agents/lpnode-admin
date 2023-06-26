package design

import (
	. "goa.design/goa/v3/dsl"
)

var LpMonit_TaskItem = Type("LpMointTaskItem", func() {
	Attribute("_id", String)
	Attribute("name", String)
	Attribute("cron", String, "定时任务")
	Attribute("createdAt", Int64, "创建时间戳")
	Attribute("scriptPath", String, "脚本路径")
	Attribute("taskType", String, "任务类型")

	Required("name", "cron", "createdAt", "taskType")
})
var _ = Service("lpmonit", func() {
	Description("监控脚本程序")
	Method("add_script", func() {
		Description("add script and save")
		Payload(func() {
			Description("Multipart request Payload")
			Attribute("name", String)
			Attribute("cron", String)
			Attribute("scriptBody", String)
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("task_id", String, "创建后的Id")
			Attribute("result", String, "创建成功后的Id")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/monit/script_add")
		})
	})
	Method("list_script", func() {
		Description("task_list")
		Payload(func() {
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(LpMonit_TaskItem), "任务列表")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/monit/task_list")
		})
	})
	Method("delete_script", func() {
		Description("task_list_delete")
		Payload(func() {
			Attribute("_id", String, "Mongodb 的主键")
			Required("_id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "是否删除成功")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/monit/task_del")
		})
	})
	Method("run_script", func() {
		Description("task_run")
		Payload(func() {
			Attribute("scriptContent", String)
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", String, "")
			Attribute("message", String, "")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/monit/script_run")
		})
	})
	Method("run_result", func() {
		Description("run_result")
		Payload(func() {
			Attribute("scriptName", String)
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", String, "")
			Attribute("message", String, "")
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/monit/run_result")
		})
	})
})

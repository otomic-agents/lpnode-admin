package design

import (
	. "goa.design/goa/v3/dsl"
)

var LpMonit_TaskItem = Type("LpMointTaskItem", func() {
	Attribute("_id", String)
	Attribute("name", String)
	Attribute("cron", String, "scheduled task")
	Attribute("createdAt", Int64, "create timestamp")
	Attribute("scriptPath", String, "script path")
	Attribute("taskType", String, "task type")

	Required("name", "cron", "createdAt", "taskType")
})
var _ = Service("lpmonit", func() {
	Description("monitor script")
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
			Attribute("task_id", String, "id after creation")
			Attribute("result", String, "id after successful creation")
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
			Attribute("result", ArrayOf(LpMonit_TaskItem), "task list")
			Attribute("message", String)
		})
		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/monit/task_list")
		})
	})
	Method("delete_script", func() {
		Description("task_list_delete")
		Payload(func() {
			Attribute("_id", String, "mongodb primary key")
			Required("_id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "result")
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

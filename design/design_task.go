package design

import (
	. "goa.design/goa/v3/dsl"
)

var Task_TaskItem = Type("TaskItem", func() {
	Attribute("_id", String)
	Attribute("schedule", String) // 定时的规则
	Attribute("taskType", String, func() {
		Enum("build-in")
		Enum("customize")
	})
	Attribute("deployed", Boolean) //  是否已经部署
	Attribute("deployMessage", String)
	Attribute("scriptPath", String)
	Attribute("scriptBody", String)
})
var Task_Deploy_CMD = Type("Task_Deploy", func() {
	Attribute("_id", String)
})
var _ = Service("taskManager", func() {
	Description("用于列表tasklist")
	Method("taskList", func() {
		Payload(func() {})
		Result(func() {
			Attribute("code", Int64)
			Attribute("result", ArrayOf(Task_TaskItem))
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/taskManager/list")
		})
	})
	Method("taskDeploy", func() {
		Payload(Task_Deploy_CMD)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "创建影响的行数")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/taskManager/deploy")
		})
	})
	Method("unDeploy", func() {
		Payload(Task_Deploy_CMD)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "创建影响的行数")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/taskManager/undeploy")
		})
	})
	Method("taskCreate", func() {
		Payload(Task_TaskItem)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Int64, "创建影响的行数")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/taskManager/create")
		})
	})
})

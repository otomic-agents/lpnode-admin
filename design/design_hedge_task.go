package design

import (
	// Import Goa DSL
	. "goa.design/goa/v3/dsl"
)

// --- Hedge Task Related Type Definitions (Prefixed with HT_) ---

// HT_HedgeTask defines the basic information of a hedge task
var HT_HedgeTask = Type("HT_HedgeTask", func() {
	Description("Basic information of a hedge task")
	Attribute("id", String, "Task ID", func() { Example("60c72b2f9b1d8e001c8a4d5e") })
	Attribute("name", String, "Task Name", func() { Example("ETH-USDT Hedge Task") })
	Attribute("amm_name", String, "AMM Name", func() { Example("Uniswap") })
	Attribute("status", String, "Task Status (pending, active, completed)", func() { Example("active") })
	Attribute("account_id", String, "Associated Account ID", func() { Example("60c72b2f9b1d8e001c8a4d5f") })
	Attribute("bridge_id", String, "Associated Bridge ID", func() { Example("60c72b2f9b1d8e001c8a4d60") })
	Attribute("created_at", String, "Creation Time", func() { Example("2023-04-20T10:00:00Z") })
	Attribute("updated_at", String, "Update Time", func() { Example("2023-04-20T11:00:00Z") })
	Attribute("created_by", String, "Creator", func() { Example("admin") })
	Required("id", "name", "status", "account_id", "bridge_id", "created_at", "updated_at")
})

// HT_CreateTaskPayload defines the request parameters for creating a hedge task
var HT_CreateTaskPayload = Type("HT_CreateTaskPayload", func() {
	Description("Request parameters for creating a hedge task")
	Attribute("id", String, "TaskId", func() { Example("60c72b2f9b1d8e001c8a4d5e") })
	Attribute("name", String, "Task Name", func() { Example("ETH-USDT Hedge Task") })
	Attribute("account_id", String, "Associated Account ID", func() { Example("60c72b2f9b1d8e001c8a4d5f") })
	Attribute("bridge_id", String, "Associated Bridge ID", func() { Example("60c72b2f9b1d8e001c8a4d60") })
	Attribute("amm_name", String, "AMM Name", func() { Example("Uniswap") })
	Required("name", "account_id", "bridge_id", "amm_name")
})

// HT_TokenBalance defines token balance information
var HT_TokenBalance = Type("HT_TokenBalance", func() {
	Description("Token balance information")
	Attribute("token", String, "Token symbol", func() { Example("USDT") })
	Attribute("tokenAddress", String, "Token contract address", func() { Example("0xdac17f958d2ee523a2206206994597c13d831ec7") })
	Attribute("cex", String, "CEX balance", func() { Example("10000") })
	Attribute("dex", String, "DEX balance", func() { Example("1000115.28820246") })
	Attribute("total", Float64, "Total balance", func() { Example(1010115.28820246) })
	Attribute("wallet", String, "Wallet address", func() { Example("0xCb4284dFA16429762e40d01F5Cff4D4bD0870f42") })
	Attribute("walletName", String, "Wallet name", func() { Example("B1") })
	Attribute("chainId", Int64, "Chain ID", func() { Example(1) })
})

// HT_InitialBalances defines the initial balances structure
var HT_InitialBalances = Type("HT_InitialBalances", func() {
	Description("Initial token balances")
	Attribute("source", HT_TokenBalance, "Source token balance")
	Attribute("destination", HT_TokenBalance, "Destination token balance")
})

// HT_RiskConfig defines the risk configuration structure
var HT_RiskConfig = Type("HT_RiskConfig", func() {
	Description("Risk configuration for hedge tasks")
	Attribute("auto_hedge", Boolean, "Auto hedge flag", func() { Example(false) })
	Attribute("max_asset_exposure", MapOf(String, Float64), "Maximum asset exposure", func() {
		Example(map[string]float64{
			"ETH":  10,
			"USDT": 10000,
			"USDC": 10000,
		})
	})
	Attribute("min_hedge_amount", MapOf(String, Float64), "Minimum hedge amount", func() {
		Example(map[string]float64{
			"ETH":  0.1,
			"USDT": 100,
			"USDC": 100,
		})
	})
	Attribute("hedge_mode", String, "Hedge mode", func() { Example("SPOT") })
})

// HT_SaveHedgeDataPayload defines the request parameters for saving hedge data
var HT_SaveHedgeDataPayload = Type("HT_SaveHedgeDataPayload", func() {
	Description("Request parameters for saving hedge data")
	Attribute("id", String, "Hedge Task ID", func() { Example("6809ffda61342a03ea15124b") })
	Attribute("hedgeTaskId", String, "Hedge Task ID (duplicate)", func() { Example("6809ffda61342a03ea15124b") })
	Attribute("initialBalances", HT_InitialBalances, "Initial balances information")
	Attribute("name", String, "Hedge name", func() { Example("0011") })
	Attribute("bridge_id", String, "Bridge ID", func() { Example("67c0175f7aecc30830cac46e") })
	Attribute("cex_account_id", String, "CEX Account ID", func() { Example("68020ad9378b05e70089f71f") })
	Attribute("amm_name", String, "AMM Name", func() { Example("amm-01") })
	Attribute("risk_config", HT_RiskConfig, "Risk configuration")
	Required("id", "initialBalances", "name", "bridge_id", "cex_account_id", "amm_name", "risk_config")
})

// HT_SaveHedgeDataResult defines the response for saving hedge data
var HT_SaveHedgeDataResult = Type("HT_SaveHedgeDataResult", func() {
	Description("Result of saving hedge data")
	Attribute("id", String, "Hedge ID", func() { Example("6809ffda61342a03ea15124b") })
	Attribute("name", String, "Hedge name", func() { Example("0011") })
	Attribute("status", String, "Status of the hedge", func() { Example("active") })
	Required("id", "name", "status")
})

// hedgeTasks service definition
var _ = Service("hedgeTasks", func() {
	Description("Manage hedge tasks")

	// listTasks method: List all hedge tasks
	Method("listTasks", func() {
		Description("List all hedge tasks")
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(HT_HedgeTask), "Hedge task list")
			Attribute("message", String, "")
			Required("code", "message", "result")
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/hedge/tasks")
			Response(StatusOK)
		})
	})
	// getTask method: Get a single hedge task by ID
	Method("getTask", func() {
		Description("Get a single hedge task by ID")
		Payload(func() {
			Attribute("task_id", String, "ID of the hedge task to retrieve", func() {
				Example("60c72b2f9b1d8e001c8a4d5e")
			})
			Required("task_id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", HT_HedgeTask, "Hedge task details")
			Attribute("message", String, "")
			Required("code", "message", "result")
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/hedge/tasks/{task_id}")
			Params(func() {
				Param("task_id", String, "Hedge task ID")
			})
			Response(StatusOK)
		})
	})
	// deleteTask method: Delete specified hedge task
	Method("deleteTask", func() {
		Description("Delete specified hedge task")
		Payload(func() {
			Attribute("task_id", String, "ID of the hedge task to delete", func() {
				Example("60c72b2f9b1d8e001c8a4d5e")
			})
			Required("task_id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("message", String, "")
			Required("code", "message")
		})

		HTTP(func() {
			DELETE("/lpnode/lpnode_admin_panel/hedge/tasks/{task_id}")
			Params(func() {
				Param("task_id", String, "Hedge task ID")
			})
			Response(StatusOK)
		})
	})

	// createTask method: Create new hedge task
	Method("createTask", func() {
		Description("Create new hedge task")
		Payload(HT_CreateTaskPayload)
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", HT_HedgeTask, "Created hedge task information")
			Attribute("message", String, "")
			Required("code", "message", "result")
		})

		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/hedge/tasks")
			Response(StatusCreated)
		})
	})

	// Add this method to the hedgeTasks service
	Method("saveHedgeData", func() {
		Description("Save hedge data to update an existing hedge configuration")
		Payload(HT_SaveHedgeDataPayload)
		Result(func() {
			Attribute("code", Int64, "Response code")
			Attribute("result", HT_SaveHedgeDataResult, "Updated hedge information")
			Attribute("message", String, "Response message")
			Required("code", "message", "result")
		})

		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/hedge/tasks/data")
			Response(StatusOK)
		})
	})
	// closeTask method: Close a hedge task
	Method("closeTask", func() {
		Description("Close a hedge task")
		Payload(func() {
			Attribute("task_id", String, "ID of the hedge task to close", func() {
				Example("60c72b2f9b1d8e001c8a4d5e")
			})
			Required("task_id")
		})
		Result(func() {
			Attribute("code", Int64, "Response code")
			Attribute("message", String, "Response message")
			Required("code", "message")
		})

		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/hedge/tasks/close/{task_id}")
			Params(func() {
				Param("task_id", String, "Hedge task ID")
			})
			Response(StatusOK)
		})
	})
})

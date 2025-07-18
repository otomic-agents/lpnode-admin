package design

import (
	. "goa.design/goa/v3/dsl"
)

// Define ChainTransaction related types
var DSG_CT_TransactionItem = Type("DSG_CT_TransactionItem", func() {
	Attribute("id", String, "Transaction ID")
	Attribute("systemChainId", String, "System chain ID")
	Attribute("chainId", Int64, "Chain ID")
	Attribute("eventName", String, "Event name")
	Attribute("lastSend", Int64, "Last send timestamp")
	Attribute("gasPrice", String, "Gas price")
	Attribute("businessId", String, "Business ID")
	Attribute("createdAt", Int64, "Creation timestamp")
	Attribute("updatedAt", Int64, "Update timestamp")
	Attribute("status", String, "Transaction status")
	Attribute("srcTokenName", String, "")
	Attribute("dstTokenName", String, "")
	Attribute("srcChainName", String, "")
	Attribute("dstChainName", String, "")

	Attribute("sendList", ArrayOf("DSG_CT_SendRecord"), "Send records")
	Attribute("gasList", ArrayOf("DSG_CT_GasRecord"), "Gas records")
})

var DSG_CT_SendRecord = Type("DSG_CT_SendRecord", func() {
	Attribute("timestamp", Int64, "Timestamp")
	Attribute("gasPrice", String, "Gas price")
	Attribute("txHash", String, "txHash")
	Attribute("transactionHash", String, "Transaction hash")
	Attribute("status", String, "Status")
	Attribute("error", String, "Error message")
	Attribute("retryCount", Int, "Retry count")
})

var DSG_CT_GasRecord = Type("DSG_CT_GasRecord", func() {
	Attribute("timestamp", Int64, "Timestamp")
	Attribute("transactionHash", String, "Transaction hash")
	Attribute("gasPrice", String, "Gas price")
})

// Define pagination result type
var DSG_CT_ListResult = Type("DSG_CT_ListResult", func() {
	Attribute("total", Int64, "Total records")
	Attribute("page", Int, "Current page")
	Attribute("pageCount", Int, "pageCount")
	Attribute("list", ArrayOf(DSG_CT_TransactionItem), "Transaction list")
	Required("total", "page", "pageCount", "list")
})

var DSG_CT_ListPayload = Type("DSG_CT_ListTransactionPayload", func() {
	Attribute("page", Int, "Page number", func() {
		Default(1)
		Minimum(1)
	})
	Attribute("pageSize", Int, "Page size", func() {
		Default(10)
		Minimum(1)
		Maximum(100)
	})
	Attribute("businessId", String, "Filter by business ID")
	Attribute("status", String, "Filter by status")
	Attribute("chainId", Int64, "Filter by chain ID")
	Attribute("startTime", Int64, "Start time")
	Attribute("endTime", Int64, "End time")
})

var _ = Service("chainClientTransaction", func() {
	Description("Blockchain transaction management service")

	// List transactions
	Method("transactionList", func() {
		Description("Get transaction list")

		Payload(DSG_CT_ListPayload)

		Result(func() {
			Attribute("code", Int64, "Status code")
			Attribute("message", String, "Response message")
			Attribute("result", DSG_CT_ListResult, "Pagination result")
			Required("code", "message", "result")
		})

		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/chain-client-transaction/list")
			// Use Body() to specify request body format
			Body(func() {
				Attribute("page")
				Attribute("pageSize")
				Attribute("businessId")
				Attribute("status")
				Attribute("chainId")
				Attribute("startTime")
				Attribute("endTime")
			})
			Response(StatusOK)
		})
	})
})

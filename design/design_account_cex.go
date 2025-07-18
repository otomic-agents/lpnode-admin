package design

import (
	// Import Goa DSL
	. "goa.design/goa/v3/dsl"
)

// --- CEX account related type definitions (with CA_ prefix) ---

// CA_CexAccountBalance defines CEX account balance information
var CA_CexAccountBalance = Type("CA_CexAccountBalance", func() {
	Description("Balance information for a single asset in a CEX account")
	Attribute("asset", String, "Asset name (e.g.: USDT, BTC)", func() { Example("USDT") })
	Attribute("total", String, "Total balance", func() { Example("1000.50") })
	Attribute("free", String, "Available balance", func() { Example("950.00") })
	Attribute("locked", String, "Locked balance", func() { Example("50.50") })
	Attribute("price", String, "Current price (optional, if provided by exchange)", func() { Example("1.00") }) // Keep original definition
	Required("asset", "total", "free", "locked")
})

// CA_CexAccount defines the core information of a CEX account (for response, without sensitive information)
var CA_CexAccount = Type("CA_CexAccount", func() {
	Description("CEX account information")
	Attribute("id", String, "ID", func() { Example("99191991") })
	Attribute("name", String, "Account name (user-defined)", func() { Example("Binance Hedge Account 1") })
	Attribute("exchange", String, "Exchange name (e.g.: binance, okx)", func() { Example("binance") })
	Attribute("api_key", String, "API Key (partially masked in list responses)", func() { Example("abc...def") })
	Attribute("status", String, "Account status (e.g.: active, inactive, invalid_keys)", func() { Example("active") })
	Required("name", "exchange", "api_key", "status")
})

// CA_CexAccountPayload defines the data required to create a CEX account
var CA_CexAccountPayload = Type("CA_CexAccountPayload", func() {
	Description("Payload for creating a CEX account")
	Attribute("name", String, "Account name (user-defined)", func() {
		Example("Binance Hedge Account 1")
	})
	Attribute("exchange", String, "Exchange name (lowercase, e.g.: binance, okx)", func() {
		Example("binance")
	})
	Attribute("api_key", String, "API Key provided by the exchange", func() {
		Example("abc123def456")
	})
	Attribute("api_secret", String, "API Secret provided by the exchange", func() {
		Example("xyz789uvw101")
	})
	Attribute("passphrase", String, "Passphrase for API key (if required by the exchange)", func() {
		Example("mysecretpassphrase") // Optional
	})

	Required("name", "exchange", "api_key", "api_secret")
})

// --- New: CA_TokenBalance defines single token balance information ---
var CA_TokenBalance = Type("CA_TokenBalance", func() {
	Description("Balance information for a single token")
	Attribute("asset", String, "Asset name", func() { Example("BTC") })
	Attribute("free", String, "Available balance", func() { Example("1") })
	Attribute("locked", String, "Locked balance", func() { Example("0") })
	Attribute("total", String, "Total balance", func() { Example("1") })
	Required("asset", "free", "locked", "total")
})

// accountCex service definition
var _ = Service("accountCex", func() {
	Description("Manage centralized exchange (CEX) accounts and their wallet information")
	// New: tokenBalance method: Get specific token balance for a specified CEX account
	Method("getAllTokenBalances", func() {
		Description("Get all token balances for a specified CEX account")
		Payload(func() {
			Attribute("account_id", String, "CEX account ID to query", func() {
				Example("68020ad9378b05e70089f71f")
			})
			Required("account_id")
		})
		Result(func() {
			Attribute("code", Int64, "Status code")
			Attribute("result", ArrayOf(CA_TokenBalance), "All token balance information")
			Attribute("message", String, "Response message")
			Required("code", "message", "result")
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/account/cex/{account_id}/tokens")
			Params(func() {
				Param("account_id", String, "CEX account ID")
			})
			Response(StatusOK)
		})
	})
	Method("tokenBalance", func() {
		Description("Get specific token balance for a specified CEX account")
		Payload(func() {
			Attribute("account_id", String, "CEX account ID to query", func() {
				Example("68020ad9378b05e70089f71f")
			})
			Attribute("symbol", String, "Token symbol to query", func() {
				Example("BTC")
			})
			Required("account_id", "symbol")
		})
		Result(func() {
			Attribute("code", Int64, "Status code")
			Attribute("result", CA_TokenBalance, "Token balance information")
			Attribute("message", String, "Response message")
			Required("code", "message", "result")
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/account/cex/{account_id}/token/{symbol}")
			Params(func() {
				Param("account_id", String, "CEX account ID")
				Param("symbol", String, "Token symbol")
			})
			Response(StatusOK)
		})
	})
	// walletInfo method: Get wallet balance for a specified CEX account
	Method("walletInfo", func() {
		Description("Get wallet balance information for a specified CEX account")
		Payload(func() {
			Attribute("account_id", String, "CEX account ID to query", func() {
				Example("60c72b2f9b1d8e001c8a4d5e")
			})
			Required("account_id")
		})
		Result(func() {
			// Define code, message, data directly in Result
			Attribute("code", Int64, "")
			Attribute("data", ArrayOf(CA_CexAccountBalance), "CEX account balance list")
			Attribute("message", String, "")
			Required("code", "message", "data") // Ensure these fields always exist
		})

		HTTP(func() {
			// Pass account_id as path parameter
			GET("/lpnode/lpnode_admin_panel/account/cex/{account_id}/walletInfo")
			Params(func() {
				Param("account_id", String, "CEX account ID")
			})
			Response(StatusOK) // Successful response 200 OK
			// Following the provided style, don't define error responses with descriptions
			// Response("not-found", StatusNotFound)
		})
	})

	// createAccount method: Create a new CEX account configuration
	Method("createAccount", func() {
		Description("Create a new CEX account configuration")
		Payload(CA_CexAccountPayload) // Use CA_CexAccountPayload as request body

		Result(func() {
			// Define code, message, data directly in Result
			Attribute("id", String, "")
			Attribute("code", Int64, "")
			Attribute("result", CA_CexAccount, "Successfully created account information") // Return created account info (without secret)
			Attribute("message", String, "")
			Required("code", "message", "result")
		})

		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/account/cex") // POST to resource collection path
			Response(StatusCreated)                        // Successful creation uses 201 Created
			// Following the provided style, don't define error responses with descriptions
			// Response("invalid-data", StatusBadRequest)
			// Response("name-taken", StatusConflict)
			// Response("exchange-not-supported", StatusBadRequest)
			// Response("invalid-keys", StatusBadRequest)
		})
	})

	// listAccounts method: List all configured CEX accounts
	Method("listAccounts", func() {
		Description("List all configured CEX accounts")
		Payload(func() {
			// Following the provided style, if query parameters are needed, define them here and map with Param
			// Attribute("page", Int, "Page number", func() { Default(1) })
			// Attribute("limit", Int, "Items per page", func() { Default(20) })
			// Attribute("exchange", String, "Filter by exchange")
		})
		Result(func() {
			// Define code, message, data directly in Result
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(CA_CexAccount), "CEX account list") // Return account list
			Attribute("message", String, "")
			Required("code", "message", "result")
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/account/cex") // GET resource collection path
			Response(StatusOK)                            // Successful response 200 OK
			// If Payload has attributes that need to be mapped to query parameters, use Param here
			// Param("page")
			// Param("limit")
			// Param("exchange")
		})
	})
	Method("deleteAccount", func() {
		Description("Delete a specified CEX account")
		Payload(func() {
			Attribute("account_id", String, "CEX account ID to delete", func() {
				Example("60c72b2f9b1d8e001c8a4d5e")
			})
			Required("account_id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("message", String, "")
			Required("code", "message")
		})

		HTTP(func() {
			DELETE("/lpnode/lpnode_admin_panel/account/cex/{account_id}")
			Params(func() {
				Param("account_id", String, "CEX account ID")
			})
			Response(StatusOK)
		})
	})
})

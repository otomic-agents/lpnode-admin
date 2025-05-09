package design

import (
	// Import Goa DSL
	. "goa.design/goa/v3/dsl"
)

// --- Simple Market Price Types ---

// TokenPrice defines the basic price information for a token
var TokenPrice = Type("TokenPrice", func() {
	Description("Basic price information for a cryptocurrency token")
	Attribute("symbol", String, "Token symbol", func() { Example("BTC") })
	Attribute("price", Float64, "Current price in USDT", func() { Example(45000.25) })
	Required("symbol", "price")
})

// PriceResponse defines the standard response structure
var PriceResponse = Type("PriceResponse", func() {
	Description("Standard response structure for price data")
	Attribute("code", Int64, "Response code")
	Attribute("message", String, "Response message")
	Attribute("result", ArrayOf(TokenPrice), "Price data for tokens")
	Required("code", "message", "result")
})

// marketPrices service definition
var _ = Service("marketPrices", func() {
	Description("Cryptocurrency market price service")

	// getAllPrices method: Get prices for all available tokens
	Method("getAllPrices", func() {
		Description("Get current USDT prices for all available tokens")
		Result(PriceResponse)

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/market/prices")
			Response(StatusOK)
		})
	})
})

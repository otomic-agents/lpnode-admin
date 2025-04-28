package design

import (
	// Import Goa DSL
	. "goa.design/goa/v3/dsl"
)

// ER_ExchangeRate defines exchange rate information
var ER_ExchangeRate = Type("ER_ExchangeRate", func() {
	Description("Exchange rate information")
	Attribute("_id", String, "ID", func() { Example("6801e24c7aecc30830418233") })
	Attribute("symbol", String, "Trading pair symbol", func() { Example("USDC/USDT") })
	Attribute("asks", ArrayOf(ArrayOf(String)), "Ask orders list", func() {
		Example([][]string{
			{"0.9996", "8046766.140174", "0"},
			{"0.9997", "175414.550295", "0"},
		})
	})
	Attribute("bids", ArrayOf(ArrayOf(String)), "Bid orders list", func() {
		Example([][]string{
			{"0.9995", "10007154.960689", "0"},
			{"0.9994", "273695.220134", "0"},
		})
	})
	Attribute("incomingTimestamp", Int64, "Received timestamp", func() { Example(1745387400090) })
	Attribute("lastUpdateId", Int64, "Last update ID", func() { Example(0) })
	Attribute("stdSymbol", String, "Standardized trading pair symbol", func() { Example("USDC/USDT") })
	Attribute("stream", String, "Data stream type", func() { Example("spot") })
	Attribute("timestamp", Int64, "Timestamp", func() { Example(1745387398704) })
	Required("symbol", "asks", "bids", "timestamp", "stdSymbol")
})

// ER_GetAllRatesResult defines the response result for getting all rates
var ER_GetAllRatesResult = Type("ER_GetAllRatesResult", func() {
	Description("Response result for getting all rates")
	Attribute("code", Int64, "Status code", func() { Example(200) })
	Attribute("result", ArrayOf(ER_ExchangeRate), "Exchange rate information list")
	Attribute("message", String, "Response message", func() { Example("success") })
	Required("code", "message", "result")
})

// exchangeRates service definition
var _ = Service("exchangeRates", func() {
	Description("Get exchange rate information")

	// getAllRates method: Get exchange rate information for all trading pairs
	Method("getAllRates", func() {
		Description("Get exchange rate information for all trading pairs")

		Result(ER_GetAllRatesResult)

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/exchange/rates")
			Response(StatusOK)
		})
	})
})

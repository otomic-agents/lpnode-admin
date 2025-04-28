package design

import (
	. "goa.design/goa/v3/dsl"
)

var ADB_TokenBalance = Type("ADB_TokenBalance", func() {
	Attribute("tokenAddress", String, "Token contract address")
	Attribute("symbol", String, "Token symbol")
	Attribute("formattedBalance", String, "Formatted balance")
	Attribute("decimals", Int, "Token decimals")
	Attribute("isNative", Boolean, "Whether it's a native token")
	Attribute("price", String, "Token price (USD)")
	Attribute("value", String, "Token value (USD)")
	Attribute("updatedAt", String, "Update time")
})

var ADB_AddressAssetGroup = Type("ADB_AddressAssetGroup", func() {
	Attribute("walletAddress", String, "Wallet address")
	Attribute("walletNames", ArrayOf(String), "List of wallet names")
	Attribute("tokens", ArrayOf(ADB_TokenBalance), "List of token balances")
	Attribute("totalValue", String, "Total value of the address (USD)")
})

var ADB_ChainAssetGroup = Type("ADB_ChainAssetGroup", func() {
	Attribute("chainId", Int, "Chain ID")
	Attribute("chainName", String, "Chain name")
	Attribute("chainType", String, "Chain type")
	Attribute("nativeToken", String, "Native token name")
	Attribute("chainLogo", String, "Chain Logo URL")
	Attribute("addressAssets", ArrayOf(ADB_AddressAssetGroup), "List of address assets")
	Attribute("totalValue", String, "Total value on chain (USD)")
})

var ADB_WalletAssetResponse = Type("ADB_WalletAssetResponse", func() {
	Attribute("totalAddresses", Int, "Number of monitored addresses")
	Attribute("lastUpdated", String, "Last update time")
	Attribute("totalValue", String, "Total asset value (USD)")
	Attribute("chainAssets", ArrayOf(ADB_ChainAssetGroup), "List of chain assets")
})

var accountDex_DexAccountBalance = Type("DexAccountBalance", func() {
	Attribute("token", String, "")
	Attribute("tokenName", String, "")
	Attribute("amount", String, "")
	Attribute("free", String, "")
	Attribute("locked", String, "")
	Attribute("price", String)
})

var accountDex_WalletInfo = Type("ADB_WalletInfo", func() {
	Attribute("_id", String, "Wallet ID")
	Attribute("walletName", String, "Wallet name")
	Attribute("address", String, "Wallet address")
	Attribute("addressLower", String, "Lowercase wallet address")
	Attribute("chainType", String, "Chain type")
	Attribute("chainId", Int, "Chain ID")
	Attribute("walletType", String, "Wallet type")
	Attribute("signServiceEndpoint", String, "Signature service endpoint")
})

var _ = Service("accountDex", func() {
	Method("walletInfo", func() {
		Payload(func() {
			Attribute("id", String, "MongoDB ID of the wallet")
			Required("id")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", accountDex_WalletInfo, "Wallet information")
			Attribute("message", String)
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/account/dex/walletInfo")
			Param("id")
		})
	})
	Method("getWalletAssets", func() {
		Payload(func() {
			Attribute("addresses", ArrayOf(String), "List of wallet addresses, returns assets for all monitored addresses if not provided")
			Attribute("currency", String, "Display currency, default is USD", func() {
				Default("USD")
				Enum("USD", "EUR", "CNY", "JPY", "GBP")
			})
			Attribute("hideZeroBalance", Boolean, "Whether to hide zero balance assets", func() {
				Default(false)
			})
			Attribute("hideSmallBalance", Boolean, "Whether to hide small balance assets", func() {
				Default(true)
			})
			Attribute("smallBalanceThreshold", String, "Small balance threshold", func() {
				Default("1.0")
			})
		})
		Result(func() {
			Attribute("code", Int, "Status code")
			Attribute("result", ADB_WalletAssetResponse, "Wallet asset data")
			Attribute("message", String, "Response message", func() {
				Default("success")
			})
		})

		HTTP(func() {
			GET("/lpnode/lpnode_admin_panel/wallet/assets")
			Param("addresses")
			Param("currency")
			Param("hideZeroBalance")
			Param("hideSmallBalance")
			Param("smallBalanceThreshold")
		})
	})
})

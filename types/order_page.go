package types

type OrderPageAssetChangeItem struct {
	Symbol string  `json:"symbol"` // 资产符号
	Amount float64 `json:"amount"` // 数量（正数表示获得，负数表示支出）
	USD    string  `json:"usd"`    // USD价值（可选）
}
type OrderPageChainTransaction struct {
	EventName   string `json:"event_name"`
	TxHash      string `json:"tx_hash"`
	ExplorerUrl string `json:"explorer_url"`
	ChainName   string `json:"chain_name"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}

type OrderPageTransactionRow struct {
	TransactionID     string                      `json:"transaction_id"`   // 交易 ID
	TransactionTime   string                      `json:"transaction_time"` // 交易时间
	Status            string                      `json:"status"`           // 状态（成功、进行中、失败）
	Type              string                      `json:"type"`             // 交易类型（跨链 Swap）
	SourceChain       string                      `json:"source_chain"`     // 源链
	TradeStatus       string                      `json:"trade_status"`
	SrcTokenAddress   string                      `json:"src_token_address"` // 源代币地址
	DstTokenAddress   string                      `json:"dst_token_address"` // 目标代币地址
	DestinationChain  string                      `json:"destination_chain"` // 目标链
	Received          []OrderPageReceivedItem     `json:"received"`          // 接收的资产
	Pay               []OrderPagePayItem          `json:"pay"`               // 支付的资产
	GasFee            []OrderPageGasFeeItem       `json:"gas_fee"`           // Gas 费用
	TotalChanges      []OrderPageAssetChangeItem  `json:"total_changes"`
	ChainTransactions []OrderPageChainTransaction `json:"chain_transactions"`
}

type OrderPageReceivedItem struct {
	Amount string `json:"amount"` // 数量
	Symbol string `json:"symbol"` // 资产符号（如 USDT、ETH）
}

type OrderPagePayItem struct {
	Amount string `json:"amount"` // 数量
	Symbol string `json:"symbol"` // 资产符号（如 B-ETH-A、BNB）

}
type OrderPageGasFeeItem struct {
	Amount string `json:"amount"` // 数量
	Symbol string `json:"symbol"` // 资产符号（如 BNB、ETH）
	USD    string `json:"usd"`    // 对应的美元价值
}

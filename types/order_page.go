package types

type OrderPageAssetChangeItem struct {
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount"`
	USD    string  `json:"usd"`
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
	TransactionID     string                      `json:"transaction_id"`
	TransactionTime   string                      `json:"transaction_time"`
	Status            string                      `json:"status"`
	Type              string                      `json:"type"`
	SourceChain       string                      `json:"source_chain"`
	TradeStatus       string                      `json:"trade_status"`
	SrcTokenAddress   string                      `json:"src_token_address"`
	DstTokenAddress   string                      `json:"dst_token_address"`
	DestinationChain  string                      `json:"destination_chain"`
	Received          []OrderPageReceivedItem     `json:"received"`
	Pay               []OrderPagePayItem          `json:"pay"`
	GasFee            []OrderPageGasFeeItem       `json:"gas_fee"`
	TotalChanges      []OrderPageAssetChangeItem  `json:"total_changes"`
	ChainTransactions []OrderPageChainTransaction `json:"chain_transactions"`
}

type OrderPageReceivedItem struct {
	Amount string `json:"amount"`
	Symbol string `json:"symbol"`
}

type OrderPagePayItem struct {
	Amount string `json:"amount"`
	Symbol string `json:"symbol"`
}
type OrderPageGasFeeItem struct {
	Amount string `json:"amount"`
	Symbol string `json:"symbol"`
	USD    string `json:"usd"`
}

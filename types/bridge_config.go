package types

type BridgeConfigLpConfigItemBridge struct {
	SrcChainId int64  `json:"src_chain_id"`
	DstChainId int64  `json:"dst_chain_id"`
	SrcToken   string `json:"src_token"`
	DstToken   string `json:"dst_token"`
}

// configLp  base struct
type BridgeConfigLpConfigItem struct {
	Bridge BridgeConfigLpConfigItemBridge `json:"bridge"`
	Wallet struct {
		Name string `json:"name"`
	} `json:"wallet"`
	LpReceiverAddress string `json:"lp_receiver_address"`
	MsmqName          string `json:"msmq_name"`
	SrcClientUri      string `json:"src_client_uri"`
	DstClientUri      string `json:"dst_client_uri"`
	RelayApiKey       string `json:"relay_api_key"`
	RelayURI          string `json:"relay_uri"`
	EnableLimiter     bool   `json:"enableLimiter"`
}

type BridgeConfigClientConfigItem struct {
	WalletName string   `json:"wallet_name"`
	PrivateKey string   `json:"private_key"`
	TokenList  []string `json:"token_list"`
}

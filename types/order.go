package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type CenterOrder struct {
	ID       primitive.ObjectID `bson:"_id"`
	ViewInfo struct {
		SrcChainName            string `json:"srcChainName"`
		SrcChainPrecision       int64  `json:"SrcChainPrecision"`
		SrcTokenName            string `json:"srcTokenName"`
		SrcTokenPrecision       int64  `json:"srcTokenPrecision"`
		DstChainName            string `json:"dstChainName"`
		DstChainNativeTokenName string `json:"DstChainNativeTokenName"`
		DstChainPrecision       int64  `json:"dstChainPrecision"`
		DstTokenName            string `json:"dstTokenName"`
		DstTokenPrecision       int64  `json:"dstTokenPrecision"`
		ReceiverWallet          string `json:"receiverWallet"`
		PaymentWallet           string `json:"paymentWallet"`
		ReceiverAmount          string `json:"receiverAmount"`
		PaymentAmount           string `json:"paymentAmount"`
		PaymentNativeAmount     string `json:"paymentNativeAmount"`
		QuoteTimestamp          int64  `json:"quoteTimestamp"`
		SrcChainRpcTx           string `json:"srcChainRpcTx"`
		DstChainRpcTx           string `json:"dstChainRpcTx"`
	} `json:"ViewInfo"`
	PreBusiness struct {
		SwapAssetInformation struct {
			BridgeName      string `bson:"bridge_name" json:"bridge_name"`
			LpIdFake        string `bson:"lp_id_fake" json:"lp_id_fake"`
			Sender          string `bson:"sender" json:"sender"`
			Amount          string `bson:"amount" json:"amount"`
			DstAddress      string `bson:"dst_address" json:"dst_address"`
			DstAmount       string `bson:"dst_amount" json:"dst_amount"`
			DstNativeAmount string `bson:"dst_native_amount" json:"dst_native_amount"`
			TimeLock        int64  `bson:"time_lock" json:"time_lock"`
			Quote           struct {
				QuoteBase struct {
					Bridge struct {
						SrcChainId int64  `bson:"src_chain_id" json:"src_chain_id"`
						DstChainId int64  `bson:"dst_chain_id" json:"dst_chain_id"`
						SrcToken   string `bson:"src_token" json:"src_token"`
						DstToken   string `bson:"dst_token" json:"dst_token"`
						BridgeName string `bson:"bridge_name" json:"bridge_name"`
					} `bson:"bridge" json:"bridge"`
					LpBridgeAddress  string `bson:"lp_bridge_address" json:"lp_bridge_address"`
					Price            string `bson:"price" json:"price"`
					NativeTokenPrice string `bson:"native_token_price" json:"native_token_price"`
					NativeTokenMax   string `bson:"native_token_max" json:"native_token_max"`
					NativeTokenMin   string `bson:"native_token_min" json:"native_token_min"`
					Capacity         string `bson:"capacity" json:"capacity"`
					LpNodeURI        string `bson:"lp_node_uri" json:"lp_node_uri"`
					QuoteHash        string `bson:"quote_hash" json:"quote_hash"`
				} `bson:"quote_base" json:"quote_base"`
				QuoteName string `bson:"quote_name" json:"quote_name"`
				Timestamp int64  `bson:"timestamp" json:"timestamp"`
			} `bson:"quote" json:"quote"`
			AppendInformation string `bson:"append_information" json:"append_information"`
			Requestor         string `bson:"requestor"`
		} `bson:"swap_asset_information" json:"swap_asset_information"`
		Hash            string `bson:"hash" json:"hash"`
		LpSalt          string `bson:"lp_salt" json:"lp_salt"`
		HashlockEvm     string `bson:"hashlock_evm" json:"hashlock_evm"`
		HashlockXrp     string `bson:"hashlock_xrp" json:"hashlock_xrp"`
		HashlockNear    string `bson:"hashlock_near" json:"hashlock_near"`
		Locked          bool   `bson:"locked" json:"locked"`
		Timestamp       int64  `bson:"timestamp" json:"timestamp"`
		OrderAppendData string `bson:"order_append_data" json:"order_append_data"`
	} `bson:"pre_business" json:"pre_business"`
	Business struct {
		BusinessID   int64  `bson:"business_id" json:"business_id"`
		Step         int64  `bson:"step" json:"step"`
		BusinessHash string `bson:"business_hash" json:"business_hash"`
	} `bson:"business" json:"business"`
	EventTransferOut struct {
		TransferOutID int64  `bson:"transfer_out_id" json:"transfer_out_id"`
		BusinessID    int64  `bson:"business_id" json:"business_id"`
		TransferInfo  string `bson:"transfer_info" json:"transfer_info"`
		TransferID    string `bson:"transfer_id" json:"transfer_id"`
		Sender        string `bson:"sender" json:"sender"`
		Receiver      string `bson:"receiver" json:"receiver"`
		Token         string `bson:"token" json:"token"`
		Amount        string `bson:"amount" json:"amount"`
		HashLock      string `bson:"hash_lock" json:"hash_lock"`
		TimeLock      int64  `bson:"time_lock" json:"time_lock"`
		DstChainID    int64  `bson:"dst_chain_id" json:"dst_chain_id"`
		DstAddress    string `bson:"dst_address" json:"dst_address"`
		BidID         string `bson:"bid_id" json:"bid_id"`
		DstToken      string `bson:"dst_token" json:"dst_token"`
		DstAmount     string `bson:"dst_amount" json:"dst_amount"`
	} `bson:"event_transfer_out" json:"event_transfer_out"`
	EventTransferIn struct {
		TransferID       string      `bson:"transfer_id" json:"transfer_id"`
		Sender           string      `bson:"sender" json:"sender"`
		Receiver         string      `bson:"receiver" json:"receiver"`
		Token            string      `bson:"token" json:"token"`
		TokenAmount      string      `bson:"token_amount" json:"token_amount"`
		EthAmount        string      `bson:"eth_amount" json:"eth_amount"`
		HashLock         string      `bson:"hash_lock" json:"hash_lock"`
		HashLockOriginal interface{} `bson:"hash_lock_original" json:"hash_lock_original"`
		TimeLock         int64       `bson:"time_lock" json:"time_lock"`
		SrcChainID       int64       `bson:"src_chain_id" json:"src_chain_id"`
		SrcTransferID    string      `bson:"src_transfer_id" json:"src_transfer_id"`
		TransferInfo     string      `bson:"transfer_info" json:"transfer_info"`
	} `bson:"event_transfer_in" json:"event_transfer_in"`
	EventTransferOutConfirm struct {
		BusinessID   int    `bson:"business_id" json:"business_id"`
		TransferInfo string `bson:"transfer_info" json:"transfer_info"`
		TransferID   string `bson:"transfer_id" json:"transfer_id"`
		Preimage     string `bson:"preimage" json:"preimage"`
	} `bson:"event_transfer_out_confirm" json:"event_transfer_out_confirm"`
	EventTransferInConfirm struct {
		BusinessID   int64  `bson:"business_id" json:"business_id"`
		TransferInfo string `bson:"transfer_info" json:"transfer_info"`
		TransferID   string `bson:"transfer_id" json:"transfer_id"`
		Preimage     string `bson:"preimage" json:"preimage"`
	} `bson:"event_transfer_in_confirm" json:"event_transfer_in_confirm"`
	EventTransferOutRefund interface{} `bson:"event_transfer_out_refund" json:"event_transfer_out_refund"`
	EventTransferInRefund  interface{} `bson:"event_transfer_in_refund" json:"event_transfer_in_refund"`
}

type AmmContext struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Summary      string             `bson:"summary" json:"summary"`
	HedgeEnabled bool               `bson:"hedgeEnabled" json:"hedgeEnabled"`
	SystemInfo   struct {
		MsmqName string `bson:"msmqName" json:"msmqName"`
	} `bson:"systemInfo" json:"systemInfo"`
	WalletInfo struct {
		WalletName string `bson:"walletName" json:"walletName"`
	} `bson:"walletInfo" json:"walletInfo"`
	AskInfo struct {
		Cid string `bson:"cid" json:"cid"`
	} `bson:"AskInfo" json:"AskInfo"`
	BaseInfo struct {
		
		SrcToken struct {
			Precision    int    `bson:"precision" json:"precision"`
			CexPrecision int    `bson:"cexPrecision" json:"cexPrecision"`
			Address      string `bson:"address" json:"address"`
			CoinType     string `bson:"coinType" json:"coinType"`
			Symbol       string `bson:"symbol" json:"symbol"`
			ChainID      int    `bson:"chainId" json:"chainId"`
		} `bson:"srcToken" json:"srcToken"`
		DstToken struct {
			Precision    int    `bson:"precision" json:"precision"`
			CexPrecision int    `bson:"cexPrecision" json:"cexPrecision"`
			Address      string `bson:"address" json:"address"`
			CoinType     string `bson:"coinType" json:"coinType"`
			Symbol       string `bson:"symbol" json:"symbol"`
			ChainID      int    `bson:"chainId" json:"chainId"`
		} `bson:"dstToken" json:"dstToken"`
	} `bson:"baseInfo" json:"baseInfo"`
	SwapInfo struct {
		InputAmount     string  `bson:"inputAmount" json:"inputAmount"`
		SrcAmount       string  `bson:"srcAmount" json:"srcAmount"`
		DstAmount       string  `bson:"dstAmount" json:"dstAmount"`
		SrcAmountNumber float64 `bson:"srcAmountNumber" json:"srcAmountNumber"`
		DstAmountNumber float64 `bson:"dstAmountNumber" json:"dstAmountNumber"`
	} `bson:"swapInfo"`
	QuoteInfo struct {
		UsdPrice         float64 `bson:"usd_price" json:"usd_price"`
		Price            string  `bson:"price" json:"price"`
		OrigPrice        string  `bson:"origPrice" json:"origPrice"`
		MinAmount        string  `bson:"min_amount" json:"min_amount"`
		Gas              string  `bson:"gas" json:"gas"`
		Capacity         string  `bson:"capacity" json:"capacity"`
		NativeTokenPrice string  `bson:"native_token_price" json:"native_token_price"`
		NativeTokenMax   string  `bson:"native_token_max" json:"native_token_max"`
		NativeTokenMin   string  `bson:"native_token_min" json:"native_token_min"`
		Timestamp        int64   `bson:"timestamp" json:"timestamp"`
		QuoteHash        string  `bson:"quote_hash" json:"quote_hash"`
		Orderbook        struct {
			A struct {
				Bids [][]float64 `bson:"bids" json:"bids"`
				Asks [][]float64 `bson:"asks" json:"asks"`
			} `bson:"A" json:"A"`
			B struct {
				Bids [][]float64 `bson:"bids" json:"bids"`
				Asks [][]float64 `bson:"asks" json:"asks"`
			} `bson:"B" json:"B"`
		} `bson:"orderbook" json:"orderbook"`
		AssetName           string  `bson:"assetName" json:"assetName"`
		AssetTokenName      string  `bson:"assetTokenName" json:"assetTokenName"`
		AssetChainInfo      string  `bson:"assetChainInfo" json:"assetChainInfo"`
		QuoteOrderbookType  string  `bson:"quote_orderbook_type" json:"quote_orderbook_type"`
		NativeTokenSymbol   string  `bson:"native_token_symbol" json:"native_token_symbol"`
		NativeTokenMinUsd   string  `bson:"native_token_min_usd" json:"native_token_min_usd"`
		NativeTokenMinCount string  `bson:"native_token_min_count" json:"native_token_min_count"`
		GasUsd              float64 `bson:"gas_usd" json:"gas_usd"`
		CapacityNum         string  `bson:"capacity_num" json:"capacity_num"`
	} `bson:"quoteInfo" json:"quoteInfo"`
	AskTime     int64 `bson:"askTime" json:"askTime"`
	SystemOrder struct {
		BridgeConfig struct {
			BridgeName string `bson:"bridge_name" json:"bridge_name"`
			SrcChainID int    `bson:"src_chain_id" json:"src_chain_id"`
			DstChainID int    `bson:"dst_chain_id" json:"dst_chain_id"`
			SrcToken   string `bson:"srcToken" json:"srcToken"`
			DstToken   string `bson:"dstToken" json:"dstToken"`
			MsmqName   string `bson:"msmq_name" json:"msmq_name"`
			Wallet     struct {
				Name    string `bson:"name" json:"name"`
				Balance struct {
				} `bson:"balance" json:"balance"`
			} `bson:"wallet" json:"wallet"`
			DstChainClientURI string `bson:"dst_chain_client_uri" json:"dst_chain_client_uri"`
		} `bson:"bridgeConfig" json:"bridgeConfig"`
		BalanceLockedID string `bson:"balanceLockedId" json:"balanceLockedId"`
		Hash            string `bson:"hash" json:"hash"`
		BaseInfo        struct {
			SrcChain struct {
				ID   int    `bson:"id" json:"id"`
				Name string `bson:"name" json:"name"`
			} `bson:"srcChain" json:"srcChain"`
			DstChain struct {
				ID   int    `bson:"id" json:"id"`
				Name string `bson:"name" json:"name"`
			} `bson:"dstChain" json:"dstChain"`
			SrcToken struct {
				Address  string `bson:"address" json:"address"`
				Symbol   string `bson:"symbol" json:"symbol"`
				CoinType string `bson:"coinType" json:"coinType"`
			} `bson:"srcToken" json:"srcToken"`
			DstToken struct {
				Address  string `bson:"address" json:"address"`
				Symbol   string `bson:"symbol" json:"symbol"`
				CoinType string `bson:"coinType" json:"coinType"`
			} `bson:"dstToken" json:"dstToken"`
		} `bson:"baseInfo" json:"baseInfo"`
		QuoteInfo struct {
			Amount           string `bson:"amount" json:"amount"`
			QuoteHash        string `bson:"quoteHash" json:"quoteHash"`
			Price            string `bson:"price" json:"price"`
			Capacity         string `bson:"capacity" json:"capacity"`
			NativeTokenPrice string `bson:"nativeTokenPrice" json:"nativeTokenPrice"`
		} `bson:"quoteInfo" json:"quoteInfo"`
		OrderID         int `bson:"orderId" json:"orderId"`
		TransferOutInfo struct {
			Amount string `bson:"amount" json:"amount"`
		} `bson:"transferOutInfo" json:"transferOutInfo"`
		TransferOutTimestamp int64 `bson:"transferOutTimestamp" json:"transferOutTimestamp"`
		CexResult            struct {
			OrderInfo struct {
				Side               string      `bson:"side" json:"side"`
				Type               string      `bson:"type" json:"type"`
				TimeInForce        string      `bson:"timeInForce" json:"timeInForce"`
				Fee                interface{} `bson:"fee" json:"fee"`
				Info               string      `bson:"info" json:"info"`
				Symbol             string      `bson:"symbol" json:"symbol"`
				StdSymbol          string      `bson:"stdSymbol" json:"stdSymbol"`
				Amount             float64     `bson:"amount" json:"amount"`
				Filled             float64     `bson:"filled" json:"filled"`
				Remaining          int         `bson:"remaining" json:"remaining"`
				ClientOrderID      string      `bson:"clientOrderId" json:"clientOrderId"`
				Timestamp          int64       `bson:"timestamp" json:"timestamp"`
				LastTradeTimestamp int64       `bson:"lastTradeTimestamp" json:"lastTradeTimestamp"`
				Average            float64     `bson:"average" json:"average"`
				AveragePrice       string      `bson:"averagePrice" json:"averagePrice"`
				Status             string      `bson:"status" json:"status"`
				FeeView            string      `bson:"feeView" json:"feeView"`
			} `bson:"orderInfo" json:"orderInfo"`
		} `bson:"cexResult" json:"cexResult"`
	} `bson:"systemOrder" json:"systemOrder"`
	LockInfo struct {
		Price string `bson:"price" json:"price"`
		Time  int64  `bson:"time" json:"time"`
	} `bson:"lockInfo" json:"lockInfo"`
}

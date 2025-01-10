package types

import (
	"time"
)

type BusinessOrder struct {
	ID                      *string                               `json:"_id" bson:"_id"`
	AppName                 *string                               `json:"app_name" bson:"appName"`
	HedgeEnabled            *bool                                 `json:"hedge_enabled" bson:"hedgeEnabled"`
	Summary                 *string                               `json:"summary" bson:"summary"`
	TradeStatus             interface{}                           `json:"tradeStatus" bson:"tradeStatus"`
	ChainOptInfo            *BusinessOrderChainOptInfo            `json:"chain_opt_info" bson:"chainOptInfo"`
	SystemInfo              *BusinessOrderSystemInfo              `json:"system_info" bson:"systemInfo"`
	WalletInfo              *BusinessOrderWalletInfo              `json:"wallet_info" bson:"walletInfo"`
	AskInfo                 *BusinessOrderAskInfo                 `json:"ask_info" bson:"AskInfo"`
	BaseInfo                *BusinessOrderBaseInfo                `json:"base_info" bson:"baseInfo"`
	SwapInfo                *BusinessOrderSwapInfo                `json:"swap_info" bson:"swapInfo"`
	QuoteInfo               *BusinessOrderQuoteInfo               `json:"quote_info" bson:"quoteInfo"`
	AskTime                 *int64                                `json:"ask_time" bson:"askTime"`
	SystemOrder             *BusinessOrderSystemOrder             `json:"system_order" bson:"systemOrder"`
	LockInfo                *BusinessOrderLockInfo                `json:"lock_info" bson:"lockInfo"`
	ProfitStatus            *int                                  `json:"profit_status" bson:"profitStatus"`
	FlowStatus              *string                               `json:"flow_status" bson:"flowStatus"`
	CreateTime              *time.Time                            `json:"create_time" bson:"createtime"`
	Version                 *int                                  `json:"__v" bson:"__v"`
	SystemContext           *BusinessOrderSystemContext           `json:"system_context" bson:"systemContext"`
	DexTradeInfoOut         *BusinessOrderDexTradeInfoOut         `json:"dex_trade_info_out" bson:"dexTradeInfo_out"`
	IsTrading               *bool                                 `json:"is_trading" bson:"isTrading"`
	DexTradeInfoIn          *BusinessOrderDexTradeInfoIn          `json:"dex_trade_info_in" bson:"dexTradeInfo_in"`
	DexTradeInfoOutConfirm  *BusinessOrderDexTradeInfoOutConfirm  `json:"dex_trade_info_out_confirm" bson:"dexTradeInfo_out_confirm"`
	DexTradeInfoInConfirm   *BusinessOrderDexTradeInfoInConfirm   `json:"dex_trade_info_in_confirm" bson:"dexTradeInfo_in_confirm"`
	DexTradeInfoInRefund    *BusinessOrderDexTradeInfoInRefund    `json:"dexTradeInfo_in_refund" bson:"dexTradeInfo_in_refund"`
	DexTradeInfoInitSwap    *BusinessOrderDexTradeInfoInitSwap    `json:"dexTradeInfo_init_swap" bson:"dexTradeInfo_init_swap"`
	DexTradeInfoConfirmSwap *BusinessOrderDexTradeInfoConfirmSwap `json:"dexTradeInfo_confirm_swap" bson:"dexTradeInfo_confirm_swap"`
	DexTradeInfoRefundSwap  *BusinessOrderDexTradeInfoRefundSwap  `json:"dexTradeInfo_refund_swap" bson:"dexTradeInfo_refund_swap"`
}
type BusinessOrderDexTradeInfoIn struct {
	RawData *BusinessOrderDexTradeInfoInRawData `json:"raw_data" bson:"rawData"`
}

type BusinessOrderDexTradeInfoInRefund struct {
	RawData *BusinessOrderDexTradeInfoInRefundRawData `json:"rawData" bson:"rawData"`
}

type BusinessOrderDexTradeInfoInRefundRawData struct {
	Class        *string `json:"@class" bson:"@class"`
	BusinessID   *int    `json:"business_id" bson:"business_id"`
	TransferInfo *string `json:"transfer_info" bson:"transfer_info"`
	TransferID   *string `json:"transfer_id" bson:"transfer_id"`
}

// BusinessOrderDexTradeInfoInitSwap 定义初始化交换的结构
type BusinessOrderDexTradeInfoInitSwap struct {
	RawData *BusinessOrderDexTradeInfoInitSwapRawData `json:"raw_data" bson:"rawData"`
}

type BusinessOrderDexTradeInfoInitSwapRawData struct {
	Class                  *string `json:"@class" bson:"@class"`
	TransferInfo           *string `json:"transfer_info" bson:"transfer_info"`
	TransferID             *string `json:"transfer_id" bson:"transfer_id"`
	Sender                 *string `json:"sender" bson:"sender"`
	Receiver               *string `json:"receiver" bson:"receiver"`
	Token                  *string `json:"token" bson:"token"`
	Amount                 *string `json:"amount" bson:"amount"`
	DstToken               *string `json:"dst_token" bson:"dst_token"`
	DstAmount              *string `json:"dst_amount" bson:"dst_amount"`
	ExpectedSingleStepTime *int    `json:"expected_single_step_time" bson:"expected_single_step_time"`
	AgreementReachedTime   *int64  `json:"agreement_reached_time" bson:"agreement_reached_time"`
	BidID                  *string `json:"bid_id" bson:"bid_id"`
	Requestor              *string `json:"requestor" bson:"requestor"`
	LPID                   *string `json:"lp_id" bson:"lp_id"`
}

// BusinessOrderDexTradeInfoConfirmSwap 定义确认交换的结构
type BusinessOrderDexTradeInfoConfirmSwap struct {
	RawData *BusinessOrderDexTradeInfoConfirmSwapRawData `json:"raw_data" bson:"rawData"`
}

type BusinessOrderDexTradeInfoConfirmSwapRawData struct {
	Class        *string `json:"@class" bson:"@class"`
	TransferInfo *string `json:"transfer_info" bson:"transfer_info"`
	TransferID   *string `json:"transfer_id" bson:"transfer_id"`
}
type BusinessOrderDexTradeInfoInRawData struct {
	Class                  *string `json:"@class" bson:"@class"`
	UUID                   *string `json:"uuid" bson:"uuid"`
	TransferID             *string `json:"transfer_id" bson:"transfer_id"`
	Sender                 *string `json:"sender" bson:"sender"`
	Receiver               *string `json:"receiver" bson:"receiver"`
	Token                  *string `json:"token" bson:"token"`
	TokenAmount            *string `json:"token_amount" bson:"token_amount"`
	EthAmount              *string `json:"eth_amount" bson:"eth_amount"`
	HashLock               *string `json:"hash_lock" bson:"hash_lock"`
	HashLockOriginal       *string `json:"hash_lock_original" bson:"hash_lock_original"`
	AgreementReachedTime   *int64  `json:"agreement_reached_time" bson:"agreement_reached_time"`
	ExpectedSingleStepTime *int    `json:"expected_single_step_time" bson:"expected_single_step_time"`
	TolerantSingleStepTime *int    `json:"tolerant_single_step_time" bson:"tolerant_single_step_time"`
	EarliestRefundTime     *int64  `json:"earliest_refund_time" bson:"earliest_refund_time"`
	SrcChainID             *int    `json:"src_chain_id" bson:"src_chain_id"`
	SrcTransferID          *string `json:"src_transfer_id" bson:"src_transfer_id"`
	TransferInfo           *string `json:"transfer_info" bson:"transfer_info"`
}
type BusinessOrderChainOptInfo struct {
	SrcChainReceiveAmount              *string  `json:"src_chain_receive_amount" bson:"srcChainReceiveAmount"`
	SrcChainReceiveAmountNumber        *float64 `json:"src_chain_receive_amount_number" bson:"srcChainReceiveAmountNumber"`
	DstChainPayAmount                  *string  `json:"dst_chain_pay_amount" bson:"dstChainPayAmount"`
	DstChainPayAmountNumber            *float64 `json:"dst_chain_pay_amount_number" bson:"dstChainPayAmountNumber"`
	DstChainPayNativeTokenAmount       *string  `json:"dst_chain_pay_native_token_amount" bson:"dstChainPayNativeTokenAmount"`
	DstChainPayNativeTokenAmountNumber *float64 `json:"dst_chain_pay_native_token_amount_number" bson:"dstChainPayNativeTokenAmountNumber"`
}

type BusinessOrderSystemInfo struct {
	MsmqName *string `json:"msmq_name" bson:"msmqName"`
}

type BusinessOrderWalletInfo struct {
	WalletName *string `json:"wallet_name" bson:"walletName"`
}

type BusinessOrderAskInfo struct {
	CID  *string `json:"cid" bson:"cid"`
	LPID *string `json:"lp_id" bson:"lpId"`
}

type BusinessOrderBaseInfo struct {
	Fee       *float64            `json:"fee" bson:"fee"`
	SrcChain  *BusinessOrderChain `json:"src_chain" bson:"srcChain"`
	DstChain  *BusinessOrderChain `json:"dst_chain" bson:"dstChain"`
	SrcToken  *BusinessOrderToken `json:"src_token" bson:"srcToken"`
	DstToken  *BusinessOrderToken `json:"dst_token" bson:"dstToken"`
	SourceFee *float64            `json:"source_fee" bson:"sourceFee"`
}

type BusinessOrderChain struct {
	ID                   *int    `json:"id" bson:"id"`
	NativeTokenName      *string `json:"native_token_name" bson:"nativeTokenName"`
	NativeTokenPrecision *int    `json:"native_token_precision" bson:"nativeTokenPrecision"`
	TokenName            *string `json:"token_name" bson:"tokenName"`
}

type BusinessOrderToken struct {
	Precision    *int    `json:"precision" bson:"precision"`
	CexPrecision *int    `json:"cex_precision" bson:"cexPrecision"`
	Address      *string `json:"address" bson:"address"`
	CoinType     *string `json:"coin_type" bson:"coinType"`
	Symbol       *string `json:"symbol" bson:"symbol"`
	ChainID      *int    `json:"chain_id" bson:"chainId"`
	TokenName    *string `json:"token_name" bson:"tokenName"`
}

type BusinessOrderSwapInfo struct {
	InputAmount           *string  `json:"input_amount" bson:"inputAmount"`
	InputAmountNumber     *float64 `json:"input_amount_number" bson:"inputAmountNumber"`
	SystemSrcFee          *float64 `json:"system_src_fee" bson:"systemSrcFee"`
	SystemDstFee          *float64 `json:"system_dst_fee" bson:"systemDstFee"`
	LPReceiveAmount       *float64 `json:"lp_receive_amount" bson:"lpReceiveAmount"`
	SrcAmount             *string  `json:"src_amount" bson:"srcAmount"`
	DstAmount             *string  `json:"dst_amount" bson:"dstAmount"`
	SrcAmountNumber       *float64 `json:"src_amount_number" bson:"srcAmountNumber"`
	DstSourceAmount       *string  `json:"dst_source_amount" bson:"dstSourceAmount"`
	DstAmountNumber       *float64 `json:"dst_amount_number" bson:"dstAmountNumber"`
	DstNativeAmount       *string  `json:"dst_native_amount" bson:"dstNativeAmount"`
	DstSourceNativeAmount *string  `json:"dst_source_native_amount" bson:"dstSourceNativeAmount"`
	DstNativeAmountNumber *float64 `json:"dst_native_amount_number" bson:"dstNativeAmountNumber"`
	StepTimeLock          *int     `json:"step_time_lock" bson:"stepTimeLock"`
}

type BusinessOrderQuoteInfo struct {
	OrigTotalPrice       *string                 `json:"orig_total_price" bson:"origTotalPrice"`
	Price                *string                 `json:"price" bson:"price"`
	OrigPrice            *string                 `json:"orig_price" bson:"origPrice"`
	DstUsdPrice          *float64                `json:"dst_usd_price" bson:"dst_usd_price"`
	MinAmount            *string                 `json:"min_amount" bson:"min_amount"`
	Gas                  *string                 `json:"gas" bson:"gas"`
	Capacity             *string                 `json:"capacity" bson:"capacity"`
	NativeTokenPrice     *string                 `json:"native_token_price" bson:"native_token_price"`
	NativeTokenUsdPrice  *string                 `json:"native_token_usdt_price" bson:"native_token_usdt_price"`
	NativeTokenMax       *string                 `json:"native_token_max" bson:"native_token_max"`
	NativeTokenMin       *string                 `json:"native_token_min" bson:"native_token_min"`
	Timestamp            *int64                  `json:"timestamp" bson:"timestamp"`
	QuoteHash            *string                 `json:"quote_hash" bson:"quote_hash"`
	QuoteOrderbookType   *string                 `json:"quote_orderbook_type" bson:"quote_orderbook_type"`
	SrcUsdPrice          *string                 `json:"src_usd_price" bson:"src_usd_price"`
	Mode                 *string                 `json:"mode" bson:"mode"`
	Orderbook            *BusinessOrderOrderbook `json:"orderbook" bson:"orderbook"`
	NativeTokenOrigPrice *string                 `json:"native_token_orig_price" bson:"native_token_orig_price"`
	NativeTokenSymbol    *string                 `json:"native_token_symbol" bson:"native_token_symbol"`
	AssetName            *string                 `json:"asset_name" bson:"assetName"`
	AssetTokenName       *string                 `json:"asset_token_name" bson:"assetTokenName"`
	AssetChainInfo       *string                 `json:"asset_chain_info" bson:"assetChainInfo"`
	CapacityNum          *string                 `json:"capacity_num" bson:"capacity_num"`
	NativeTokenMinNumber *float64                `json:"native_token_min_number" bson:"native_token_min_number"`
	NativeTokenMaxNumber *float64                `json:"native_token_max_number" bson:"native_token_max_number"`
	InputAmount          *string                 `json:"input_amount" bson:"inputAmount"`
}

type BusinessOrderOrderbook struct {
	A *interface{}             `json:"A" bson:"A"`
	B *BusinessOrderOrderbookB `json:"B" bson:"B"`
}

type BusinessOrderOrderbookB struct {
	Bids      *[][]float64 `json:"bids" bson:"bids"`
	Asks      *[][]float64 `json:"asks" bson:"asks"`
	Timestamp *int64       `json:"timestamp" bson:"timestamp"`
}
type BusinessOrderDexTradeInfoRefundSwap struct {
	RawData *BusinessOrderDexTradeInfoRefundSwapRawData `json:"raw_data" bson:"rawData"`
}
type BusinessOrderDexTradeInfoRefundSwapRawData struct {
	Class        *string `json:"@class" bson:"@class"`
	TransferInfo *string `json:"transfer_info" bson:"transfer_info"`
	TransferID   *string `json:"transfer_id" bson:"transfer_id"`
}
type BusinessOrderSystemOrder struct {
	BridgeConfig                *BusinessOrderBridgeConfig    `json:"bridge_config" bson:"bridgeConfig"`
	BalanceLockedID             *string                       `json:"balance_locked_id" bson:"balanceLockedId"`
	Hash                        *string                       `json:"hash" bson:"hash"`
	BaseInfo                    *BusinessOrderBaseInfo        `json:"base_info" bson:"baseInfo"`
	QuoteInfo                   *BusinessOrderQuoteInfo       `json:"quote_info" bson:"quoteInfo"`
	OrderID                     *int                          `json:"order_id" bson:"orderId"`
	TransferOutInfo             *BusinessOrderTransferOutInfo `json:"transfer_out_info" bson:"transferOutInfo"`
	TransferOutTimestamp        *int64                        `json:"transfer_out_timestamp" bson:"transferOutTimestamp"`
	TransferOutConfirmTimestamp *int64                        `json:"transfer_out_confirm_timestamp" bson:"transferOutConfirmTimestamp"`
	TransferInTimestamp         *int64                        `json:"transfer_in_timestamp" bson:"transferInTimestamp"`
	TransferInConfirmTimestamp  *int64                        `json:"transfer_in_confirm_timestamp" bson:"transferInConfirmTimestamp"`
	TransferInRefundTimestamp   *int64                        `json:"transfer_in_refund_timestamp" bson:"transferInRefundTimestamp"`
	InitSwapTimestamp           *int64                        `json:"init_swap_timestamp" bson:"initSwapTimestamp"`
	ConfirmSwapTimestamp        *int64                        `json:"confirm_swap_timestamp" bson:"confirmSwapTimestamp"`
	RefundSwapTimestamp         *int64                        `json:"refund_swap_timestamp" bson:"refundSwapTimestamp"`
	HedgeResult                 []map[string]interface{}      `json:"hedge_result" bson:"hedgeResult"`
}

type BusinessOrderBridgeConfig struct {
	ID                *string              `json:"id" bson:"id"`
	BridgeName        *string              `json:"bridge_name" bson:"bridge_name"`
	SrcChainID        *int                 `json:"src_chain_id" bson:"src_chain_id"`
	DstChainID        *int                 `json:"dst_chain_id" bson:"dst_chain_id"`
	SrcToken          *string              `json:"src_token" bson:"srcToken"`
	DstToken          *string              `json:"dst_token" bson:"dstToken"`
	MsmqName          *string              `json:"msmq_name" bson:"msmq_name"`
	MsmqPath          *string              `json:"msmq_path" bson:"msmq_path"`
	Wallet            *BusinessOrderWallet `json:"wallet" bson:"wallet"`
	Fee               *string              `json:"fee" bson:"fee"`
	DstChainClientURI *string              `json:"dst_chain_client_uri" bson:"dst_chain_client_uri"`
	SrcChainClientURL *string              `json:"src_chain_client_url" bson:"src_chain_client_url"`
	EnableHedge       *bool                `json:"enable_hedge" bson:"enable_hedge"`
	RelayAPIKey       *string              `json:"relay_api_key" bson:"relay_api_key"`
	FeeManager        *interface{}         `json:"fee_manager____" bson:"fee_manager____"`
	LPWalletInfo      *interface{}         `json:"lp_wallet_info____" bson:"lp_wallet_info____"`
}

type BusinessOrderWallet struct {
	Name    *string      `json:"name" bson:"name"`
	Balance *interface{} `json:"balance" bson:"balance"`
}

type BusinessOrderTransferOutInfo struct {
	Amount *string `json:"amount" bson:"amount"`
}

type BusinessOrderLockInfo struct {
	Fee              *string  `json:"fee" bson:"fee"`
	Price            *string  `json:"price" bson:"price"`
	NativeTokenPrice *string  `json:"native_token_price" bson:"nativeTokenPrice"`
	Time             *int64   `json:"time" bson:"time"`
	DstTokenPrice    *float64 `json:"dst_token_price" bson:"dstTokenPrice"`
	SrcTokenPrice    *string  `json:"src_token_price" bson:"srcTokenPrice"`
}

type BusinessOrderSystemContext struct {
	LockStepInfo *BusinessOrderLockStepInfo `json:"lock_step_info" bson:"lockStepInfo"`
}

type BusinessOrderLockStepInfo struct {
	Class           *string                   `json:"@class" bson:"@class"`
	Cmd             *string                   `json:"cmd" bson:"cmd"`
	QuoteData       *interface{}              `json:"quote_data" bson:"quote_data"`
	QuoteRemoveInfo *interface{}              `json:"quote_remove_info" bson:"quote_remove_info"`
	PreBusiness     *BusinessOrderPreBusiness `json:"pre_business" bson:"pre_business"`
}

type BusinessOrderPreBusiness struct {
	Class                *string                            `json:"@class" bson:"@class"`
	SwapAssetInformation *BusinessOrderSwapAssetInformation `json:"swap_asset_information" bson:"swap_asset_information"`
	Hash                 *string                            `json:"hash" bson:"hash"`
	LPSalt               *interface{}                       `json:"lp_salt" bson:"lp_salt"`
	HashlockEVM          *interface{}                       `json:"hashlock_evm" bson:"hashlock_evm"`
	HashlockXRP          *interface{}                       `json:"hashlock_xrp" bson:"hashlock_xrp"`
	HashlockNear         *interface{}                       `json:"hashlock_near" bson:"hashlock_near"`
	HashlockSolana       *interface{}                       `json:"hashlock_solana" bson:"hashlock_solana"`
	Locked               *bool                              `json:"locked" bson:"locked"`
	LockMessage          *string                            `json:"lock_message" bson:"lock_message"`
	Timestamp            *interface{}                       `json:"timestamp" bson:"timestamp"`
	OrderAppendData      *string                            `json:"order_append_data" bson:"order_append_data"`
	IsKYC                *interface{}                       `json:"is_kyc" bson:"is_kyc"`
	SameDID              *interface{}                       `json:"same_did" bson:"same_did"`
	KYCInfo              *BusinessOrderKYCInfo              `json:"kyc_info" bson:"kyc_info"`
}

type BusinessOrderKYCInfo struct {
	Class                *string `json:"@class" bson:"@class"`
	Address              *string `json:"address" bson:"address"`
	Birthday             *string `json:"birthday" bson:"birthday"`
	Country              *string `json:"country" bson:"country"`
	Email                *string `json:"email" bson:"email"`
	FirstName            *string `json:"first_name" bson:"first_name"`
	Gender               *string `json:"gender" bson:"gender"`
	IDEndImage           *string `json:"id_end_image" bson:"id_end_image"`
	IDFrontImage         *string `json:"id_front_image" bson:"id_front_image"`
	IDNumber             *string `json:"id_number" bson:"id_number"`
	IDType               *string `json:"id_type" bson:"id_type"`
	Image1               *string `json:"image1" bson:"image1"`
	Image2               *string `json:"image2" bson:"image2"`
	LastName             *string `json:"last_name" bson:"last_name"`
	Phone                *string `json:"phone" bson:"phone"`
	Username             *string `json:"username" bson:"username"`
	IdentificationPhoto1 *string `json:"identification_photo_1" bson:"identification_photo_1"`
	IdentificationPhoto2 *string `json:"identification_photo_2" bson:"identification_photo_2"`
	IdentificationPhoto3 *string `json:"identification_photo_3" bson:"identification_photo_3"`
	IdentificationPhoto4 *string `json:"identification_photo_4" bson:"identification_photo_4"`
	DID                  *string `json:"did" bson:"did"`
	Status               *string `json:"status" bson:"status"`
	VP                   *string `json:"vp" bson:"vp"`
}

type BusinessOrderSwapAssetInformation struct {
	BridgeName             *string             `json:"bridge_name" bson:"bridge_name"`
	LPIDFake               *string             `json:"lp_id_fake" bson:"lp_id_fake"`
	Sender                 *string             `json:"sender" bson:"sender"`
	Amount                 *string             `json:"amount" bson:"amount"`
	DstAddress             *string             `json:"dst_address" bson:"dst_address"`
	DstAmount              *string             `json:"dst_amount" bson:"dst_amount"`
	DstNativeAmount        *string             `json:"dst_native_amount" bson:"dst_native_amount"`
	SystemFeeSrc           *int                `json:"system_fee_src" bson:"system_fee_src"`
	SystemFeeDst           *int                `json:"system_fee_dst" bson:"system_fee_dst"`
	DstAmountNeed          *string             `json:"dst_amount_need" bson:"dst_amount_need"`
	DstNativeAmountNeed    *string             `json:"dst_native_amount_need" bson:"dst_native_amount_need"`
	AgreementReachedTime   *int64              `json:"agreement_reached_time" bson:"agreement_reached_time"`
	Quote                  *BusinessOrderQuote `json:"quote" bson:"quote"`
	AppendInformation      *string             `json:"append_information" bson:"append_information"`
	JWS                    *interface{}        `json:"jws" bson:"jws"`
	DID                    *interface{}        `json:"did" bson:"did"`
	Requestor              *string             `json:"requestor" bson:"requestor"`
	UserSign               *string             `json:"user_sign" bson:"user_sign"`
	LPSign                 *interface{}        `json:"lp_sign" bson:"lp_sign"`
	ExpectedSingleStepTime *int                `json:"expected_single_step_time" bson:"expected_single_step_time"`
	TolerantSingleStepTime *int                `json:"tolerant_single_step_time" bson:"tolerant_single_step_time"`
	EarliestRefundTime     *int64              `json:"earliest_refund_time" bson:"earliest_refund_time"`
	LPSignAddress          *interface{}        `json:"lp_sign_address" bson:"lp_sign_address"`
	SwapType               *string             `json:"swap_type" bson:"swap_type"`
}

type BusinessOrderDexTradeInfoOut struct {
	RawData *BusinessOrderDexTradeInfoOutRawData `json:"raw_data" bson:"rawData"`
}
type BusinessOrderDexTradeInfoOutRawData struct {
	Class                  *string `json:"@class" bson:"@class"`
	UUID                   *string `json:"uuid" bson:"uuid"`
	TransferOutID          *int    `json:"transfer_out_id" bson:"transfer_out_id"`
	BusinessID             *int    `json:"business_id" bson:"business_id"`
	TransferInfo           *string `json:"transfer_info" bson:"transfer_info"`
	TransferID             *string `json:"transfer_id" bson:"transfer_id"`
	Sender                 *string `json:"sender" bson:"sender"`
	Receiver               *string `json:"receiver" bson:"receiver"`
	Token                  *string `json:"token" bson:"token"`
	Amount                 *string `json:"amount" bson:"amount"`
	HashLock               *string `json:"hash_lock" bson:"hash_lock"`
	AgreementReachedTime   *int64  `json:"agreement_reached_time" bson:"agreement_reached_time"`
	ExpectedSingleStepTime *int    `json:"expected_single_step_time" bson:"expected_single_step_time"`
	TolerantSingleStepTime *int    `json:"tolerant_single_step_time" bson:"tolerant_single_step_time"`
	EarliestRefundTime     *int64  `json:"earliest_refund_time" bson:"earliest_refund_time"`
	DstChainID             *int    `json:"dst_chain_id" bson:"dst_chain_id"`
	DstAddress             *string `json:"dst_address" bson:"dst_address"`
	BidID                  *string `json:"bid_id" bson:"bid_id"`
	DstToken               *string `json:"dst_token" bson:"dst_token"`
	DstAmount              *string `json:"dst_amount" bson:"dst_amount"`
}

type BusinessOrderDexTradeInfoOutConfirm struct {
	RawData *BusinessOrderDexTradeInfoOutConfirmRawData `json:"raw_data" bson:"rawData"`
}
type BusinessOrderDexTradeInfoOutConfirmRawData struct {
	Class        *string `json:"@class" bson:"@class"`
	BusinessID   *int    `json:"business_id" bson:"business_id"`
	TransferInfo *string `json:"transfer_info" bson:"transfer_info"`
	TransferID   *string `json:"transfer_id" bson:"transfer_id"`
	Preimage     *string `json:"preimage" bson:"preimage"`
}

type BusinessOrderDexTradeInfoInConfirm struct {
	RawData *BusinessOrderDexTradeInfoInConfirmRawData `json:"raw_data" bson:"rawData"`
}

type BusinessOrderDexTradeInfoInConfirmRawData struct {
	Class        *string `json:"@class" bson:"@class"`
	BusinessID   *int    `json:"business_id" bson:"business_id"`
	TransferInfo *string `json:"transfer_info" bson:"transfer_info"`
	TransferID   *string `json:"transfer_id" bson:"transfer_id"`
	Preimage     *string `json:"preimage" bson:"preimage"`
}

type BusinessOrderQuote struct {
	Class                 *string                             `json:"@class" bson:"@class"`
	QuoteBase             *BusinessOrderQuoteBase             `json:"quote_base" bson:"quote_base"`
	AuthenticationLimiter *BusinessOrderAuthenticationLimiter `json:"authentication_limiter" bson:"authentication_limiter"`
	QuoteName             *string                             `json:"quote_name" bson:"quote_name"`
	Timestamp             *int64                              `json:"timestamp" bson:"timestamp"`
}
type BusinessOrderQuoteBase struct {
	Class            *string              `json:"@class" bson:"@class"`
	Bridge           *BusinessOrderBridge `json:"bridge" bson:"bridge"`
	LPBridgeAddress  *string              `json:"lp_bridge_address" bson:"lp_bridge_address"`
	Price            *string              `json:"price" bson:"price"`
	NativeTokenPrice *string              `json:"native_token_price" bson:"native_token_price"`
	NativeTokenMax   *string              `json:"native_token_max" bson:"native_token_max"`
	NativeTokenMin   *string              `json:"native_token_min" bson:"native_token_min"`
	Capacity         *string              `json:"capacity" bson:"capacity"`
	LPNodeURI        *string              `json:"lp_node_uri" bson:"lp_node_uri"`
	QuoteHash        *string              `json:"quote_hash" bson:"quote_hash"`
	RelayAPIKey      *string              `json:"relay_api_key" bson:"relay_api_key"`
}

type BusinessOrderBridge struct {
	Class      *string `json:"@class" bson:"@class"`
	SrcChainID *int    `json:"src_chain_id" bson:"src_chain_id"`
	DstChainID *int    `json:"dst_chain_id" bson:"dst_chain_id"`
	SrcToken   *string `json:"src_token" bson:"src_token"`
	DstToken   *string `json:"dst_token" bson:"dst_token"`
	BridgeName *string `json:"bridge_name" bson:"bridge_name"`
}

type BusinessOrderAuthenticationLimiter struct {
	Class            *string      `json:"@class" bson:"@class"`
	CountryWhiteList *interface{} `json:"country_white_list" bson:"country_white_list"`
	CountryBlackList *interface{} `json:"country_black_list" bson:"country_black_list"`
	MinAge           *interface{} `json:"min_age" bson:"min_age"`
	AccountBlackList *interface{} `json:"account_black_list" bson:"account_black_list"`
	LimiterState     *string      `json:"limiter_state" bson:"limiter_state"`
}

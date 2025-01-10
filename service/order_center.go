package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/redis_database"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Chain ID to name mapping
var chainMap = map[int]string{
	9006: "BSC",
	501:  "Solana",
	60:   "Ethereum",
	614:  "Optimism",
}

type AssetChange struct {
	Symbol string
	Amount float64
}

// ConvertToString converts any numeric or string value to string
func ConvertToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case *string:
		return lo.FromPtrOr(v, "")
	case int:
		return strconv.FormatInt(int64(v), 10)
	case *int:
		return strconv.FormatInt(int64(lo.FromPtrOr(v, 0)), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case *int32:
		return strconv.FormatInt(int64(lo.FromPtrOr(v, 0)), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case *int64:
		return strconv.FormatInt(lo.FromPtrOr(v, 0), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case *uint:
		return strconv.FormatUint(uint64(lo.FromPtrOr(v, 0)), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case *uint32:
		return strconv.FormatUint(uint64(lo.FromPtrOr(v, 0)), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case *uint64:
		return strconv.FormatUint(lo.FromPtrOr(v, 0), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case *float32:
		return strconv.FormatFloat(float64(lo.FromPtrOr(v, 0.0)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case *float64:
		return strconv.FormatFloat(lo.FromPtrOr(v, 0.0), 'f', -1, 64)
	default:
		return ""
	}
}

type OrderCenterLogicService struct {
}

func (ocls *OrderCenterLogicService) AllRedis() (ret []types.CenterOrder, err error) {
	ret = make([]types.CenterOrder, 0)
	strList, err := redis_database.GetDataRedis().Smembers("KEY_BUSINESS_STATUS_INBUSINESS")
	if err != nil {
		return
	}
	for _, v := range strList {
		row := types.CenterOrder{}
		unmarshalErr := json.Unmarshal([]byte(v), &row)
		if unmarshalErr != nil {
			log.Println(unmarshalErr)
			continue
		}
		ret = append(ret, row)
	}
	ret, renderErr := ocls.OrderRender(ret)
	if renderErr != nil {
		err = errors.WithMessage(renderErr, "render order error:")
		return
	}
	return
}
func (ocls *OrderCenterLogicService) GetTransactionList(queryOption struct {
	Page     int64
	PageSize int64
	Status   int64
}, finder bson.M) ([]types.OrderPageTransactionRow, int64, error) {

	skip := queryOption.Page*queryOption.PageSize - queryOption.PageSize
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(skip).SetLimit(queryOption.PageSize)

	err, cursor := database.FindAllOpt("main", "ammContext_amm-01", finder, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var businessOrders []types.BusinessOrder
	if err := cursor.All(context.Background(), &businessOrders); err != nil {
		return nil, 0, err
	}

	count, err := database.Count("main", "ammContext_amm-01", finder)
	if err != nil {
		return nil, 0, err
	}
	pageCount := count / queryOption.PageSize
	if count%queryOption.PageSize != 0 {
		pageCount++
	}

	var transactionRows []types.OrderPageTransactionRow
	for _, order := range businessOrders {
		// spew.Dump(order)
		row := ocls.convertToOrderPageTransactionRow(order)
		transactionRows = append(transactionRows, row)
	}

	return transactionRows, pageCount, nil
}
func (ocls *OrderCenterLogicService) convertToOrderPageTransactionRow(order types.BusinessOrder) types.OrderPageTransactionRow {
	// Default values
	const (
		defaultTime   = "2006/01/02 15:04"
		defaultStatus = "Unknown"
		defaultSymbol = "Unknown"
	)

	// Helper function: Format amount as a string, handling precision
	formatAmount := func(amount float64, precision int) string {
		if precision <= 0 {
			precision = 8 // Default precision
		}
		// Use fmt.Sprintf to format the number and remove trailing zeros
		formatted := strings.TrimRight(strings.TrimRight(fmt.Sprintf(fmt.Sprintf("%%.%df", precision), amount), "0"), ".")
		if formatted == "" {
			return "0"
		}
		return formatted
	}

	// Helper function: Convert to decimal
	convertToDecimal := func(amount string, precision int) string {
		if amount == "" || amount == "0" {
			return "0"
		}
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			fmt.Printf("💢 Failed to parse amount: %v\n", err)
			return "0"
		}
		return formatAmount(amountFloat/math.Pow10(precision), precision)
	}

	// Convert transaction time
	transactionTime := lo.TernaryF(
		order.AskTime != nil,
		func() string { return time.Unix(*order.AskTime/1000, 0).Format("2006/01/02 15:04") },
		func() string { return time.Unix(0, 0).Format(defaultTime) },
	)

	status := ocls.convertStatus(order)
	// Build received array
	received := []types.OrderPageReceivedItem{
		{
			Amount: convertToDecimal(
				lo.FromPtrOr(order.SwapInfo.SrcAmount, ""),
				lo.FromPtrOr(order.BaseInfo.SrcToken.Precision, 0),
			),
			Symbol: lo.FromPtrOr(order.BaseInfo.SrcToken.Symbol, defaultSymbol),
		},
	}

	// Build pay array, filtering out items with zero amount
	var pay []types.OrderPagePayItem
	dstAmount := convertToDecimal(
		lo.FromPtrOr(order.SwapInfo.DstAmount, ""),
		lo.FromPtrOr(order.BaseInfo.DstToken.Precision, 0),
	)
	if dstAmount != "0" {
		pay = append(pay, types.OrderPagePayItem{
			Amount: dstAmount,
			Symbol: lo.FromPtrOr(order.BaseInfo.DstToken.Symbol, defaultSymbol),
		})
	}

	dstNativeAmount := convertToDecimal(
		lo.FromPtrOr(order.SwapInfo.DstNativeAmount, ""),
		lo.FromPtrOr(order.BaseInfo.DstChain.NativeTokenPrecision, 0),
	)
	if dstNativeAmount != "0" {
		pay = append(pay, types.OrderPagePayItem{
			Amount: dstNativeAmount,
			Symbol: lo.FromPtrOr(order.BaseInfo.DstChain.TokenName, defaultSymbol),
		})
	}

	// Convert Gas fee
	nativeTokenPrice, err := strconv.ParseFloat(lo.FromPtrOr(order.QuoteInfo.NativeTokenUsdPrice, "0"), 64)
	if err != nil {
		fmt.Printf("💢 Failed to parse NativeTokenPrice: %v\n", err)
		nativeTokenPrice = 0
	}

	fmt.Printf("🔗 Processing chain ID: %v\n", ptr.ToInt(order.BaseInfo.DstChain.ID))
	dexChanges := calculateDEXAssetChanges(order)
	cexChanges := calculateCEXAssetChanges(order.SystemOrder.HedgeResult)

	gasFee := ocls.processGasFees(order, nativeTokenPrice)
	SwapType := ptr.ToString(order.SystemContext.LockStepInfo.PreBusiness.SwapAssetInformation.SwapType)

	srcChainID := lo.FromPtrOr(order.BaseInfo.SrcChain.ID, 0)
	dstChainID := lo.FromPtrOr(order.BaseInfo.DstChain.ID, 0)

	sourceChain := "Unknown"
	if srcChainID != 0 {
		if chainName, exists := chainMap[srcChainID]; exists {
			sourceChain = chainName
		}
	}

	destinationChain := "Unknown"
	if dstChainID != 0 {
		if chainName, exists := chainMap[dstChainID]; exists {
			destinationChain = chainName
		}
	}
	if status != "Success" {
		dexChanges = []AssetChange{}
	}
	totalChanges := calculateTotalAssetChanges(dexChanges, cexChanges, gasFee)
	chainTxs := ocls.getChainTransactions(order)
	return types.OrderPageTransactionRow{
		TransactionID:     lo.FromPtrOr(order.ID, "Unknown"),
		TransactionTime:   transactionTime,
		Status:            status,
		Type:              SwapType,
		SourceChain:       sourceChain,
		TradeStatus:       ConvertToString(order.TradeStatus),
		SrcTokenAddress:   ptr.ToString(order.BaseInfo.SrcToken.Address),
		DstTokenAddress:   ptr.ToString(order.BaseInfo.DstToken.Address),
		TotalChanges:      convertToAssetChangeItems(totalChanges),
		DestinationChain:  destinationChain,
		Received:          received,
		Pay:               pay,
		GasFee:            gasFee,
		ChainTransactions: chainTxs,
	}
}
func (ocls *OrderCenterLogicService) getChainTransactions(order types.BusinessOrder) []types.OrderPageChainTransaction {
	txs := make([]types.OrderPageChainTransaction, 0)

	// Handle transfer out transaction
	if order.DexTradeInfoOut != nil && order.DexTradeInfoOut.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoOut.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "TransferOut",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.SrcChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.SrcChain.ID, 0)],
				Status:      transferInfo.Get("status").String(),
				Timestamp:   lo.FromPtrOr(order.SystemOrder.TransferOutTimestamp, lo.FromPtrOr(order.AskTime, 0)),
			})
		}
	}

	// Handle transfer out confirm transaction
	if order.DexTradeInfoOutConfirm != nil && order.DexTradeInfoOutConfirm.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoOutConfirm.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "TransferOutConfirm",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.SrcChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.SrcChain.ID, 0)],
				Status:      transferInfo.Get("status").String(),
				Timestamp:   lo.FromPtrOr(order.SystemOrder.TransferOutConfirmTimestamp, lo.FromPtrOr(order.SystemOrder.TransferOutTimestamp, lo.FromPtrOr(order.AskTime, 0))),
			})
		}
	}

	// Handle transfer in transaction
	if order.DexTradeInfoIn != nil && order.DexTradeInfoIn.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoIn.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "TransferIn",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.DstChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.DstChain.ID, 0)],
				Status:      transferInfo.Get("status").String(),
				Timestamp:   lo.FromPtrOr(order.SystemOrder.TransferInTimestamp, lo.FromPtrOr(order.SystemOrder.TransferOutConfirmTimestamp, lo.FromPtrOr(order.AskTime, 0))),
			})
		}
	}

	// Handle transfer in confirm transaction
	if order.DexTradeInfoInConfirm != nil && order.DexTradeInfoInConfirm.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoInConfirm.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "TransferInConfirm",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.DstChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.DstChain.ID, 0)],
				Status:      transferInfo.Get("status").String(),
				Timestamp:   lo.FromPtrOr(order.SystemOrder.TransferInConfirmTimestamp, lo.FromPtrOr(order.SystemOrder.TransferInTimestamp, lo.FromPtrOr(order.AskTime, 0))),
			})
		}
	}

	// Handle transfer in refund transaction
	if order.DexTradeInfoInRefund != nil && order.DexTradeInfoInRefund.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoInRefund.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "TransferInRefund",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.DstChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.DstChain.ID, 0)],
				Status:      transferInfo.Get("status").String(),
				Timestamp:   lo.FromPtrOr(order.SystemOrder.TransferInRefundTimestamp, lo.FromPtrOr(order.SystemOrder.TransferInTimestamp, lo.FromPtrOr(order.AskTime, 0))),
			})
		}
	}

	// Handle init swap transaction
	if order.DexTradeInfoInitSwap != nil && order.DexTradeInfoInitSwap.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoInitSwap.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "InitSwap",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.SrcChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.SrcChain.ID, 0)],
				Status:      "1",
				Timestamp:   lo.FromPtrOr(order.SystemOrder.InitSwapTimestamp, lo.FromPtrOr(order.DexTradeInfoInitSwap.RawData.AgreementReachedTime, 0)*1000),
			})
		}
	}

	// Handle confirm swap transaction
	if order.DexTradeInfoConfirmSwap != nil && order.DexTradeInfoConfirmSwap.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoConfirmSwap.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "ConfirmSwap",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.SrcChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.SrcChain.ID, 0)],
				Status:      transferInfo.Get("status").String(),
				Timestamp:   lo.FromPtrOr(order.SystemOrder.ConfirmSwapTimestamp, lo.FromPtrOr(order.SystemOrder.InitSwapTimestamp, lo.FromPtrOr(order.AskTime, 0))),
			})
		}
	}

	// Handle refund swap transaction
	if order.DexTradeInfoRefundSwap != nil && order.DexTradeInfoRefundSwap.RawData != nil {
		transferInfo := gjson.Parse(lo.FromPtrOr(order.DexTradeInfoRefundSwap.RawData.TransferInfo, ""))
		if txHash := transferInfo.Get("transactionHash").String(); txHash != "" {
			txs = append(txs, types.OrderPageChainTransaction{
				EventName:   "RefundSwap",
				TxHash:      txHash,
				ExplorerUrl: getExplorerUrl(order.BaseInfo.SrcChain.ID, txHash),
				ChainName:   chainMap[lo.FromPtrOr(order.BaseInfo.SrcChain.ID, 0)],
				Status:      transferInfo.Get("status").String(),
				Timestamp:   lo.FromPtrOr(order.SystemOrder.RefundSwapTimestamp, lo.FromPtrOr(order.SystemOrder.InitSwapTimestamp, lo.FromPtrOr(order.AskTime, 0))),
			})
		}
	}

	return txs
}

func getExplorerUrl(chainID *int, txHash string) string {
	if chainID == nil {
		return ""
	}

	explorerUrls := map[int]string{
		9006: "https://bscscan.com/tx/",
		501:  "https://solscan.io/tx/",
		60:   "https://etherscan.io/tx/",
		614:  "https://optimistic.etherscan.io/tx/",
	}

	if baseUrl, exists := explorerUrls[*chainID]; exists {
		return baseUrl + txHash
	}
	return ""
}
func convertToAssetChangeItems(changes []AssetChange) []types.OrderPageAssetChangeItem {
	items := make([]types.OrderPageAssetChangeItem, len(changes))
	for i, change := range changes {
		items[i] = types.OrderPageAssetChangeItem{
			Symbol: change.Symbol,
			Amount: change.Amount,
			USD:    "",
		}
	}
	return items
}
func (ocls *OrderCenterLogicService) convertStatus(order types.BusinessOrder) string {
	const (
		defaultStatus = "Unknown"
		successStatus = "Success"
		failedStatus  = "Failed"
		inProgress    = "In Progress"
	)

	if order.DexTradeInfoOut != nil {
		log.Printf("DexTradeInfoOut exists, checking transfer status...")

		if order.DexTradeInfoInConfirm != nil {
			log.Printf("DexTradeInfoInConfirm exists, returning success")
			return successStatus
		}
		if order.DexTradeInfoInRefund != nil {
			log.Printf("DexTradeInfoInRefund exists, returning failed")
			return failedStatus
		}
		log.Printf("Only DexTradeInfoOut exists, returning in progress")
		return inProgress
	}

	if order.DexTradeInfoInitSwap != nil {
		log.Printf("DexTradeInfoInitSwap exists, checking swap status...")

		if order.DexTradeInfoConfirmSwap != nil {
			if order.DexTradeInfoConfirmSwap.RawData != nil &&
				order.DexTradeInfoConfirmSwap.RawData.TransferInfo != nil {
				log.Printf("DexTradeInfoConfirmSwap exists with valid transfer_info, returning success")
				return successStatus
			}
			log.Printf("DexTradeInfoConfirmSwap exists but transfer_info is nil, returning failed")
			return failedStatus
		}

		if order.DexTradeInfoRefundSwap != nil {
			log.Printf("DexTradeInfoRefundSwap exists, returning failed")
			return failedStatus
		}

		log.Printf("Only InitSwap exists, returning in progress")
		return inProgress
	}

	log.Printf("No trade info found, returning default status")
	return defaultStatus
}
func calculateTotalAssetChanges(dexChanges []AssetChange, cexChanges []AssetChange, gasFees []types.OrderPageGasFeeItem) []AssetChange {
	fmt.Println("\n🔍 Debug Info:")
	fmt.Printf("DEX Changes: %+v\n", dexChanges)
	fmt.Printf("CEX Changes: %+v\n", cexChanges)
	fmt.Printf("Gas Fees: %+v\n", gasFees)
	// Use map to merge changes of same assets
	assetMap := make(map[string]float64)

	// Process DEX changes
	for _, change := range dexChanges {
		assetMap[change.Symbol] += change.Amount
	}

	// Process CEX changes
	for _, change := range cexChanges {
		// Remove CEX_ prefix for merging
		symbol := strings.TrimPrefix(change.Symbol, "CEX_")
		assetMap[symbol] += change.Amount
	}

	// Process Gas fees
	for _, fee := range gasFees {
		amount, _ := strconv.ParseFloat(fee.Amount, 64)
		assetMap[fee.Symbol] -= amount // Subtract gas fees
	}

	// Convert back to AssetChange slice
	var totalChanges []AssetChange
	for symbol, amount := range assetMap {
		totalChanges = append(totalChanges, AssetChange{
			Symbol: symbol,
			Amount: amount,
		})
	}

	// Print summary
	fmt.Println("\n📊 Total Asset Changes Summary (including Gas fees):")
	for _, change := range totalChanges {
		if change.Amount > 0 {
			fmt.Printf("%s +%.8f ", change.Symbol, change.Amount)
		} else {
			fmt.Printf("%s %.8f ", change.Symbol, change.Amount)
		}
	}
	fmt.Println()

	return totalChanges
}

func calculateDEXAssetChanges(data interface{}) []AssetChange {
	var changes []AssetChange

	// Convert to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("⚠️ Failed to marshal data: %v\n", err)
		return nil
	}
	jsonData := string(jsonBytes)

	// Get token information
	srcSymbol := gjson.Get(jsonData, "base_info.src_token.symbol").String()
	if srcSymbol == "" {
		fmt.Println("⚠️ Failed to get source token symbol")
		return nil
	}

	srcPrecision := gjson.Get(jsonData, "base_info.src_token.precision").Int()
	if srcPrecision == 0 {
		fmt.Println("⚠️ Failed to get source token precision")
		return nil
	}

	dstSymbol := gjson.Get(jsonData, "base_info.dst_token.symbol").String()
	if dstSymbol == "" {
		fmt.Println("⚠️ Failed to get destination token symbol")
		return nil
	}

	dstPrecision := gjson.Get(jsonData, "base_info.dst_token.precision").Int()
	if dstPrecision == 0 {
		fmt.Println("⚠️ Failed to get destination token precision")
		return nil
	}

	// Get amount information
	srcAmount := gjson.Get(jsonData, "swap_info.src_amount").String()
	if srcAmount == "" {
		fmt.Println("⚠️ Failed to get source amount")
		return nil
	}

	dstAmount := gjson.Get(jsonData, "swap_info.dst_amount").String()
	if dstAmount == "" {
		fmt.Println("⚠️ Failed to get destination amount")
		return nil
	}

	fmt.Println("\n🔄 Processing DEX Swap:")
	fmt.Printf("📥 Input: %s (%s)\n", srcSymbol, srcAmount)
	fmt.Printf("📤 Output: %s (%s)\n", dstSymbol, dstAmount)

	// Calculate source token changes
	srcAmountFloat := utils.ConvertWeiToFloat(srcAmount, int(srcPrecision))
	changes = append(changes, AssetChange{
		Symbol: srcSymbol,
		Amount: srcAmountFloat,
	})

	// Calculate destination token changes
	dstAmountFloat := utils.ConvertWeiToFloat(dstAmount, int(dstPrecision))
	changes = append(changes, AssetChange{
		Symbol: dstSymbol,
		Amount: -dstAmountFloat,
	})

	// Process native token consumption
	dstNativeAmount := gjson.Get(jsonData, "swap_info.dst_native_amount").String()
	if dstNativeAmount != "0" && dstNativeAmount != "" {
		nativeTokenName := gjson.Get(jsonData, "base_info.dst_chain.native_token_name").String()
		nativeAmountFloat := utils.ConvertWeiToFloat(dstNativeAmount, 18)
		changes = append(changes, AssetChange{
			Symbol: nativeTokenName,
			Amount: -nativeAmountFloat,
		})
		fmt.Printf("💰 Gas Fee: %s (%s)\n", nativeTokenName, dstNativeAmount)
	}

	// Print asset changes summary
	fmt.Println("\n📊 DEX Asset Changes Summary:")
	for _, change := range changes {
		if change.Amount > 0 {
			fmt.Printf("%s +%.8f ", change.Symbol, change.Amount)
		} else {
			fmt.Printf("%s %.8f ", change.Symbol, change.Amount)
		}
	}
	fmt.Println("\n✅ DEX swap processed successfully")

	return changes
}

func calculateCEXAssetChanges(hedgeResults []map[string]interface{}) []AssetChange {
	if len(hedgeResults) == 0 {
		fmt.Println("⚠️ No hedge results to process")
		return nil
	}

	cexAssetChanges := make(map[string]float64)
	processedOrders := 0

	for i, result := range hedgeResults {
		// Get plan and execution result
		plan, ok := result["plan"].(map[string]interface{})
		if !ok {
			fmt.Printf("⚠️ Failed to get plan data for result %d\n", i+1)
			continue
		}
		execResult, ok := result["result"].(map[string]interface{})
		if !ok {
			fmt.Printf("⚠️ Failed to get execution result data for result %d\n", i+1)
			continue
		}

		// Get trading pair information
		symbol, ok := plan["symbol"].(string)
		if !ok {
			fmt.Printf("⚠️ Invalid symbol format in result %d\n", i+1)
			continue
		}
		symbols := strings.Split(symbol, "/")
		if len(symbols) != 2 {
			fmt.Printf("⚠️ Invalid trading pair format: %s in result %d\n", symbol, i+1)
			continue
		}
		baseAsset := symbols[0]  // ETH
		quoteAsset := symbols[1] // USDT

		// Get trade direction and amount
		side, ok := execResult["side"].(string)
		if !ok {
			fmt.Printf("⚠️ Failed to get trade side in result %d\n", i+1)
			continue
		}

		amount, ok := execResult["amount"].(float64)
		if !ok {
			fmt.Printf("⚠️ Failed to get amount in result %d\n", i+1)
			continue
		}

		price, ok := execResult["average"].(string)
		if !ok {
			fmt.Printf("⚠️ Failed to get price in result %d\n", i+1)
			continue
		}

		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Printf("⚠️ Failed to parse price '%s' in result %d: %v\n", price, i+1, err)
			continue
		}

		// Calculate the total cost
		cost := amount * priceFloat

		// Process fees
		if fee, ok := execResult["fee"].(map[string]interface{}); ok {
			for feeCurrency, feeAmount := range fee {
				feeAmountStr, ok := feeAmount.(string)
				if !ok {
					fmt.Printf("⚠️ Invalid fee format for %s in result %d\n", feeCurrency, i+1)
					continue
				}
				feeAmountFloat, err := strconv.ParseFloat(feeAmountStr, 64)
				if err != nil {
					fmt.Printf("⚠️ Failed to parse fee amount '%s' in result %d: %v\n", feeAmountStr, i+1, err)
					continue
				}
				cexAssetChanges[feeCurrency] -= feeAmountFloat
			}
		}

		// Update asset changes based on trade direction
		switch strings.ToUpper(side) {
		case "BUY":
			cexAssetChanges[baseAsset] += amount // Buy base asset
			cexAssetChanges[quoteAsset] -= cost  // Pay with quote asset
		case "SELL":
			cexAssetChanges[baseAsset] -= amount // Sell base asset
			cexAssetChanges[quoteAsset] += cost  // Receive quote asset
		default:
			fmt.Printf("⚠️ Unknown trade side '%s' in result %d\n", side, i+1)
			continue
		}

		processedOrders++
	}

	if processedOrders == 0 {
		fmt.Println("⚠️ No orders were successfully processed")
		return nil
	}

	// Convert map to slice of AssetChange structs
	var result []AssetChange
	for asset, change := range cexAssetChanges {
		result = append(result, AssetChange{
			Symbol: "CEX_" + asset,
			Amount: change,
		})
	}

	// Print summary of asset changes
	fmt.Println("\n📊 CEX Asset Changes Summary:")
	for _, change := range result {
		if change.Amount > 0 {
			fmt.Printf("%s +%.8f ", change.Symbol, change.Amount)
		} else {
			fmt.Printf("%s %.8f ", change.Symbol, change.Amount)
		}
	}
	fmt.Printf("\n✅ Successfully processed %d out of %d orders\n", processedOrders, len(hedgeResults))

	return result
}

func (ocls *OrderCenterLogicService) processGasFees(order types.BusinessOrder, nativeTokenPrice float64) []types.OrderPageGasFeeItem {
	var gasFee []types.OrderPageGasFeeItem

	// Process dexTradeInfo_in
	if order.DexTradeInfoIn != nil && order.DexTradeInfoIn.RawData != nil {
		transferInfo := lo.FromPtrOr(order.DexTradeInfoIn.RawData.TransferInfo, "")
		// fmt.Printf("📝 dexTradeInfo_in transferInfo: %s\n", transferInfo)
		if transferInfo != "" {
			if item := ocls.extractGasFeeByChain(order, transferInfo, nativeTokenPrice); item.Amount != "0" {
				gasFee = append(gasFee, item)
			}
		}
	}

	// Process dexTradeInfo_in_confirm
	if order.DexTradeInfoInConfirm != nil && order.DexTradeInfoInConfirm.RawData != nil {
		transferInfo := lo.FromPtrOr(order.DexTradeInfoInConfirm.RawData.TransferInfo, "")
		// fmt.Printf("📝 dexTradeInfo_in_confirm transferInfo: %s\n", transferInfo)
		if transferInfo != "" {
			if item := ocls.extractGasFeeByChain(order, transferInfo, nativeTokenPrice); item.Amount != "0" {
				gasFee = append(gasFee, item)
			}
		}
	}
	// Process dexTradeInfo_in_refund
	if order.DexTradeInfoInRefund != nil && order.DexTradeInfoInRefund.RawData != nil {
		transferInfo := lo.FromPtrOr(order.DexTradeInfoInRefund.RawData.TransferInfo, "")
		// fmt.Printf("📝 dexTradeInfo_in_refund transferInfo: %s\n", transferInfo)
		if transferInfo != "" {
			if item := ocls.extractGasFeeByChain(order, transferInfo, nativeTokenPrice); item.Amount != "0" {
				gasFee = append(gasFee, item)
			}
		}
	}
	// Process dexTradeInfo_confirm_swap
	if order.DexTradeInfoConfirmSwap != nil && order.DexTradeInfoConfirmSwap.RawData != nil {
		transferInfo := lo.FromPtrOr(order.DexTradeInfoConfirmSwap.RawData.TransferInfo, "")
		// fmt.Printf("📝 dexTradeInfo_confirm_swap transferInfo: %s\n", transferInfo)
		if transferInfo != "" {
			if item := ocls.extractGasFeeByChain(order, transferInfo, nativeTokenPrice); item.Amount != "0" {
				gasFee = append(gasFee, item)
			}
		}
	}

	return gasFee
}

func (ocls *OrderCenterLogicService) extractGasFeeByChain(order types.BusinessOrder, transferInfo string, nativeTokenPrice float64) types.OrderPageGasFeeItem {
	// fmt.Printf("🔍 Attempting to parse transferInfo: %s\n", transferInfo)

	if !json.Valid([]byte(transferInfo)) {
		fmt.Printf("💢 Invalid JSON in transferInfo: %s\n", transferInfo)
		return types.OrderPageGasFeeItem{
			Amount: "0",
			Symbol: "",
			USD:    "0",
		}
	}

	if ptr.ToInt(order.BaseInfo.DstChain.ID) == 501 {
		return extractSolanaGasFee(transferInfo, nativeTokenPrice)
	}
	return extractGasFee(transferInfo,
		lo.FromPtrOr(order.BaseInfo.SrcChain.NativeTokenPrecision, 18),
		nativeTokenPrice)
}

func extractGasFee(transferInfo string, nativeTokenPrecision int, nativeTokenUsdtPrice float64) types.OrderPageGasFeeItem {
	fmt.Println("nativeTokenUsdtPrice", nativeTokenUsdtPrice)
	// Default values
	defaultAmount := "0.00000000"
	defaultSymbol := "BNB"

	// Convert hex string to decimal value
	hexToDecimal := func(hexStr string) float64 {
		if hexStr == "" {
			fmt.Printf("Hex string is empty, using default value: 0\n")
			return 0
		}
		// Remove the 0x prefix
		if len(hexStr) > 2 && hexStr[:2] == "0x" {
			hexStr = hexStr[2:]
		}
		// Convert hex string to decimal
		value, err := strconv.ParseUint(hexStr, 16, 64)
		if err != nil {
			fmt.Printf("Failed to parse hex string '%s': %v, using default value: 0\n", hexStr, err)
			return 0
		}
		return float64(value)
	}

	// Unit conversion function
	convertToDecimal := func(amount string, precision int) float64 {
		if amount == "" {
			fmt.Printf("Amount is empty, using default value: 0\n")
			return 0
		}
		// If it's a hex string, convert it to decimal first
		if len(amount) > 2 && amount[:2] == "0x" {
			return hexToDecimal(amount) / math.Pow10(precision)
		}
		// If it's a decimal string, convert directly
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			fmt.Printf("Failed to parse amount '%s': %v, using default value: 0\n", amount, err)
			return 0
		}
		return amountFloat / math.Pow10(precision)
	}

	// Parse transferInfo
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(transferInfo), &data); err != nil {
		fmt.Printf("Failed to unmarshal transferInfo: %v\n", err)
		return types.OrderPageGasFeeItem{
			Amount: defaultAmount,
			Symbol: defaultSymbol,
			USD:    defaultAmount,
		}
	}

	// Get gasUsed and effectiveGasPrice
	gasUsed, ok := data["gasUsed"].(string)
	if !ok {
		fmt.Printf("gasUsed is missing or not a string in transferInfo\n")
		return types.OrderPageGasFeeItem{
			Amount: defaultAmount,
			Symbol: defaultSymbol,
			USD:    defaultAmount,
		}
	}
	effectiveGasPrice, ok := data["effectiveGasPrice"].(string)
	if !ok {
		fmt.Printf("effectiveGasPrice is missing or not a string in transferInfo\n")
		return types.OrderPageGasFeeItem{
			Amount: defaultAmount,
			Symbol: defaultSymbol,
			USD:    defaultAmount,
		}
	}

	// Convert gasUsed and effectiveGasPrice to numeric values
	gasUsedValue := convertToDecimal(gasUsed, 0)                                        // gasUsed is an integer, no precision conversion needed
	effectiveGasPriceValue := convertToDecimal(effectiveGasPrice, nativeTokenPrecision) // effectiveGasPrice is in wei, convert to BNB

	// Calculate actual Gas fee
	amount := gasUsedValue * effectiveGasPriceValue
	usd := amount * nativeTokenUsdtPrice

	// Format to string, keep 8 decimal places
	amountStr := fmt.Sprintf("%.8f", amount)
	usdStr := fmt.Sprintf("%.8f", usd)

	fmt.Printf("Raw Values: gasUsed=%s, effectiveGasPrice=%s\n", gasUsed, effectiveGasPrice)
	fmt.Printf("Converted Values: gasUsed=%f, effectiveGasPrice=%f (wei), effectiveGasPrice=%f (BNB)\n", gasUsedValue, hexToDecimal(effectiveGasPrice), effectiveGasPriceValue)
	fmt.Printf("Calculated Gas Fee: amount=%s BNB, usd=%s\n", amountStr, usdStr)

	return types.OrderPageGasFeeItem{
		Amount: amountStr,
		Symbol: defaultSymbol,
		USD:    usdStr,
	}
}

// Extract and calculate Solana Gas fee
func extractSolanaGasFee(transferInfo string, solanaUsdPrice float64) types.OrderPageGasFeeItem {
	defaultSymbol := "SOL"

	// Parse JSON string
	var txInfo map[string]interface{}
	if err := json.Unmarshal([]byte(transferInfo), &txInfo); err != nil {
		// If parsing fails, return default values, using string format
		return types.OrderPageGasFeeItem{
			Amount: "0.00000000",
			Symbol: defaultSymbol,
			USD:    "0.00000000",
		}
	}

	// Get gasUsed from JSON
	var gasUsed float64
	switch v := txInfo["gasUsed"].(type) {
	case float64:
		gasUsed = v
	case string:
		gasUsed, _ = strconv.ParseFloat(v, 64)
	default:
		// Return 0 if unable to parse
		return types.OrderPageGasFeeItem{
			Amount: "0.00000000",
			Symbol: defaultSymbol,
			USD:    "0.00000000",
		}
	}

	// Directly divide gasUsed by 1e9 to get SOL amount
	amount := gasUsed / 1e9
	usd := amount * solanaUsdPrice

	return types.OrderPageGasFeeItem{
		Amount: fmt.Sprintf("%.8f", amount),
		Symbol: defaultSymbol,
		USD:    fmt.Sprintf("%.8f", usd),
	}
}

func (ocls *OrderCenterLogicService) All(queryOption struct {
	Page     int64
	PageSize int64
	Status   int64
}, finder bson.M) (ret []types.CenterOrder, pageCount int64, err error) {
	var results []types.CenterOrder
	skip := queryOption.Page*queryOption.PageSize - queryOption.PageSize
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(skip).SetLimit(queryOption.PageSize)
	pageCount = 0

	count, err := database.Count("businessHistory", "business", finder)
	if err != nil {
		return
	}
	err, cursor := database.FindAllOpt("businessHistory", "business", finder, opts)
	pageCount = count / queryOption.PageSize
	if count%queryOption.PageSize != 0 {
		pageCount++
	}
	if pageCount == 0 {
		pageCount = 1
	}
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "cursor all error")
		return
	}

	for _, result := range results {
		// log.Println(cursor.Current.String())
		cursor.Decode(&result)
	}

	ret, renderErr := ocls.OrderRender(results)
	if renderErr != nil {
		err = errors.WithMessage(renderErr, "render order error:")
		return
	}
	return
}
func (ocls *OrderCenterLogicService) OrderRender(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	log.Println("OrderRender")
	var bridgeNameList []string = []string{}
	seenBridge := make(map[string]bool)
	ret = source
	var tokenList []string
	seenToken := make(map[string]bool)

	for _, v := range source {
		srcToken := strings.ToLower(v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken)
		dstToken := strings.ToLower(v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstToken)
		srcChainId := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcChainId
		dstChainId := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstChainId
		msmqName := fmt.Sprintf("%s/%s_%d_%d", srcToken, dstToken, srcChainId, dstChainId)
		if !seenBridge[msmqName] {
			bridgeNameList = append(bridgeNameList, msmqName)
			seenBridge[msmqName] = true
		}

		if !seenToken[srcToken] {
			tokenList = append(tokenList, srcToken)
			seenToken[srcToken] = true
		}
		if !seenToken[dstToken] {
			tokenList = append(tokenList, dstToken)
			seenToken[dstToken] = true
		}
	}
	log.Println(bridgeNameList)
	log.Println(tokenList)
	ret, renderTokenErr := ocls.orderRenderTokens(ret, tokenList)
	if renderTokenErr != nil {
		err = errors.WithMessage(renderTokenErr, "render token err:")
		return
	}

	ret, renderChainInfoErr := ocls.orderRenderChainInfo(ret)

	if renderChainInfoErr != nil {
		err = errors.WithMessage(renderChainInfoErr, "render chainInfo err:")
		return
	}

	ret, renderBridgeErr := ocls.orderRenderBridges(ret)
	if renderBridgeErr != nil {
		err = errors.WithMessage(renderBridgeErr, "render bridges err:")
		return
	}
	ret, renderViewErr := ocls.orderRenderViews(ret)
	if renderViewErr != nil {
		err = errors.WithMessage(renderViewErr, "render view err:")
		return
	}
	// pp.Print(ret)
	return
}
func (ocls *OrderCenterLogicService) orderRenderTokens(source []types.CenterOrder, tokenAddress []string) (ret []types.CenterOrder, err error) {
	log.Println(tokenAddress, "🚛🚛🚛🚛")
	ret = source
	var results []types.DBTokenRow
	err, cursor := database.FindAll("main", "tokens", bson.M{})
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	log.Println(results)
	tokenMap := make(map[string]types.DBTokenRow)
	for _, result := range results {
		if result.ChainId == 501 {
			tokenHexAddress, convertErr := utils.Base58ToHexString(result.Address)
			if convertErr != nil {
				log.Println("convert error")
				break
			}
			tokenMap[tokenHexAddress] = result
		} else {
			tokenMap[result.Address] = result
		}
	}

	for i, order := range source {
		if token, ok := tokenMap[order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", token.TokenName, token.Precision)
			order.ViewInfo.SrcTokenName = token.TokenName
			order.ViewInfo.SrcTokenPrecision = token.Precision
			source[i] = order
		} else {
			log.Println("not Found", order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken)
		}
		if token, ok := tokenMap[order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstToken]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", token.TokenName, token.Precision)
			order.ViewInfo.DstTokenName = token.TokenName
			order.ViewInfo.DstTokenPrecision = token.Precision
			source[i] = order
		}
	}
	// log.Println(tokenMap)
	for _, view := range source {
		log.Println(view.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken)
		log.Println("🎰🎰🎰🎰🎰", view.ViewInfo.SrcTokenName, view.ViewInfo.DstTokenName)
	}
	return
}
func (ocls *OrderCenterLogicService) orderRenderBridges(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	bridgeMap := make(map[string]types.DBBridgeWalletDetailsAggregateItem)
	ret = source
	mdb, err := database.GetSession("main")
	if err != nil {
		return
	}
	var results []types.DBBridgeWalletDetailsAggregateItem

	cursor, err := mdb.Collection("bridges").Aggregate(context.TODO(), bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "wallets",
				"localField":   "src_wallet_id",
				"foreignField": "_id",
				"as":           "wallet_details",
			},
		},
		bson.M{
			"$unwind": "$wallet_details",
		},
		bson.M{
			"$project": bson.M{
				"walletDetailsName":       "$wallet_details.walletName",
				"walletDetailsPrivateKey": "$wallet_details.privateKey",
				"_id":                     1,
				"wallet_id":               1,
				"msmqName":                1,
				"walletName":              1,
			},
		},
	})
	if err != nil {
		err = errors.WithMessage(err, "Aggregate err:")
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		bridgeMap[result.MsmqName] = result
	}
	for i, v := range source {
		srcToken := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken
		dstToken := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstToken
		srcChain := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcChainId
		dstChain := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstChainId
		msmqName := fmt.Sprintf("%s/%s_%d_%d", srcToken, dstToken, srcChain, dstChain)
		if bridge, ok := bridgeMap[msmqName]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁", bridge.WalletDetailsName)
			v.ViewInfo.ReceiverWallet = bridge.WalletDetailsName
			v.ViewInfo.PaymentWallet = bridge.WalletName
			source[i] = v
		}
	}

	return
}
func (ocls *OrderCenterLogicService) orderRenderChainInfo(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	ret = source
	var results []types.DbChainListRow
	err, cursor := database.FindAll("main", "chainList", bson.M{})
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	chainMap := make(map[int]types.DbChainListRow)
	for _, result := range results {
		chainMap[result.ChainId] = result
	}
	for i, order := range source {
		if chain, ok := chainMap[int(order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcChainId)]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", order.EventTransferIn.SrcChainID)
			order.ViewInfo.SrcChainName = chain.ChainName
			order.ViewInfo.SrcChainPrecision = int64(chain.Precision)
			order.ViewInfo.SrcChainRpcTx = chain.RpcTx
		}
		if chain, ok := chainMap[int(order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstChainId)]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", order.EventTransferIn.SrcChainID)
			order.ViewInfo.DstChainName = chain.ChainName
			order.ViewInfo.DstChainNativeTokenName = chain.TokenName
			order.ViewInfo.DstChainPrecision = int64(chain.Precision)
			order.ViewInfo.DstChainRpcTx = chain.RpcTx
		}
		source[i] = order
	}
	return
}
func (ocls *OrderCenterLogicService) orderRenderViews(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	for i, v := range source {
		inputAmount := new(big.Float)
		inputAmount.SetString(v.PreBusiness.SwapAssetInformation.Amount)
		srcPrecision := int(v.ViewInfo.SrcTokenPrecision)
		v.ViewInfo.ReceiverAmount = new(big.Float).Quo(inputAmount, big.NewFloat(math.Pow10(srcPrecision))).String()
		log.Println(v.PreBusiness.SwapAssetInformation.Amount, srcPrecision, v.ViewInfo.ReceiverAmount)

		sendAmount := new(big.Float)
		sendAmount.SetString(v.PreBusiness.SwapAssetInformation.DstAmount)
		dstPrecision := int(v.ViewInfo.DstTokenPrecision)
		v.ViewInfo.PaymentAmount = new(big.Float).Quo(sendAmount, big.NewFloat(math.Pow10(dstPrecision))).String()
		nativePrecision := int(v.ViewInfo.DstChainPrecision)
		sendNativeAmount := new(big.Float)
		sendNativeAmount.SetString(v.PreBusiness.SwapAssetInformation.DstNativeAmount)
		v.ViewInfo.PaymentNativeAmount = new(big.Float).Quo(sendNativeAmount, big.NewFloat(math.Pow10(nativePrecision))).String()

		// time info
		v.ViewInfo.QuoteTimestamp = v.PreBusiness.SwapAssetInformation.Quote.Timestamp

		// assign
		source[i] = v
	}
	ret = source
	return
}

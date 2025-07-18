package adminapiservice

import (
	chainclienttransaction "admin-panel/gen/chain_client_transaction"
	"admin-panel/service"
	"admin-panel/types"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"go.mongodb.org/mongo-driver/bson"
)

// chainClientTransaction service example implementation.
// The example methods log the requests and return zero values.
type chainClientTransactionsrvc struct {
	logger *log.Logger
}

// NewChainClientTransaction returns the chainClientTransaction service
// implementation.
func NewChainClientTransaction(logger *log.Logger) chainclienttransaction.Service {
	return &chainClientTransactionsrvc{logger}
}

// Get transaction list
func (s *chainClientTransactionsrvc) TransactionList(ctx context.Context, p *chainclienttransaction.DSGCTListTransactionPayload) (res *chainclienttransaction.TransactionListResult, err error) {
	res = &chainclienttransaction.TransactionListResult{}
	resList := make([]*chainclienttransaction.DSGCTTransactionItem, 0)

	queryOption := struct {
		Page     int64
		PageSize int64
		Status   int64
	}{
		Page:     1,
		PageSize: 20,
		Status:   0,
	}
	if p.Page > 0 {
		queryOption.Page = int64(p.Page)
	}
	if p.PageSize > 0 {
		queryOption.PageSize = int64(p.PageSize)
	}
	var pageCount int64 = 1
	var queryError error = nil
	var viewChainTransactions []types.View_ChainTransaction = make([]types.View_ChainTransaction, 0)
	cctx := &service.ChainClientTransactionService{}
	viewChainTransactions, pageCount, queryError = cctx.All(queryOption, "amm-01", bson.M{})
	log.Println(pageCount)
	log.Println(queryError)
	for _, viewTxRow := range viewChainTransactions {
		viewTx := viewTxRow
		log.Println(viewTx.BusinessID, viewTx.EventName)
		rowItem := &chainclienttransaction.DSGCTTransactionItem{}
		rowItem.SendList = make([]*chainclienttransaction.DSGCTSendRecord, 0)
		rowItem.BusinessID = &viewTx.BusinessID
		rowItem.SrcTokenName = &viewTx.SrcToken
		rowItem.DstTokenName = ptr.String(viewTx.DstToken)
		rowItem.EventName = ptr.String(viewTx.EventName)
		rowItem.SystemChainID = ptr.String(viewTx.SystemChainId)
		rowItem.SrcChainName = ptr.String(viewTx.SrcChainName)
		rowItem.DstChainName = ptr.String(viewTx.DstChainName)
		rowItem.GasPrice = ptr.String(viewTx.LatestGasPrice)
		rowItem.LastSend = ptr.Int64(int64(viewTx.LatestSend))
		rowItem.Status = ptr.String(viewTx.Status)

		for _, sendRecord := range viewTx.SendList {
			recordCopy := sendRecord
			// spew.Dump(&sendRecord.Status)
			rowItem.SendList = append(rowItem.SendList, &chainclienttransaction.DSGCTSendRecord{
				TxHash:    &recordCopy.TxHash,
				Timestamp: &recordCopy.Timestamp,
				GasPrice:  &recordCopy.GasPrice,
				Status:    &recordCopy.Status,
				Error:     &recordCopy.Error,
			})
		}
		for _, gasRecord := range viewTx.IncGasList {
			gasRecordCopy := gasRecord
			rowItem.GasList = append(rowItem.GasList, &chainclienttransaction.DSGCTGasRecord{
				GasPrice:  &gasRecordCopy.NewGasPrice,
				Timestamp: &gasRecordCopy.Timestamp,
			})
		}
		resList = append(resList, rowItem)

	}
	res.Code = 0
	res.Message = ""
	res.Result = &chainclienttransaction.DSGCTListResult{
		Page:      p.Page,
		PageCount: int(pageCount),
		List:      resList,
	}
	s.logger.Print("chainClientTransaction.transactionList")
	return
}

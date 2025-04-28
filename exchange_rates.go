package adminapiservice

import (
	apps_exchange_rate "admin-panel/apps/ExchangeRate"
	exchangerates "admin-panel/gen/exchange_rates"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
)

// exchangeRates service example implementation.
// The example methods log the requests and return zero values.
type exchangeRatessrvc struct {
	logger *log.Logger
}

// NewExchangeRates returns the exchangeRates service implementation.
func NewExchangeRates(logger *log.Logger) exchangerates.Service {
	return &exchangeRatessrvc{logger}
}

func convertStringArrays(arrays [][]string) [][]string {
	result := make([][]string, len(arrays))
	for i, arr := range arrays {
		result[i] = make([]string, len(arr))
		copy(result[i], arr)
	}
	return result
}

// Get exchange rate information for all trading pairs
func (s *exchangeRatessrvc) GetAllRates(ctx context.Context) (res *exchangerates.ERGetAllRatesResult, err error) {
	s.logger.Print("exchangeRates.getAllRates")

	// Initialize result
	res = &exchangerates.ERGetAllRatesResult{
		Code:    0,
		Message: "Success",
		Result:  []*exchangerates.ERExchangeRate{},
	}
	// Create and call ExchangeRateService's GetAllRates method to get exchange rate data
	erService := apps_exchange_rate.NewExchangeRateService(s.logger)
	rates, err := erService.GetAllRates(ctx)
	if err != nil {
		s.logger.Printf("Error getting exchange rates: %v", err)
		res.Code = 500
		res.Message = "Failed to retrieve exchange rates: " + err.Error()
		return res, nil
	}

	// Convert DBExchangeRate objects to ERExchangeRate objects
	for _, rate := range rates {
		erRate := &exchangerates.ERExchangeRate{
			ID:                ptr.String(rate.ID.Hex()),
			Symbol:            rate.Symbol,
			Asks:              convertStringArrays(rate.Asks),
			Bids:              convertStringArrays(rate.Bids),
			IncomingTimestamp: ptr.Int64(rate.IncomingTimestamp),

			StdSymbol: rate.StdSymbol,
			Stream:    ptr.String(rate.Stream),
			Timestamp: rate.Timestamp,
		}

		res.Result = append(res.Result, erRate)
	}

	return res, nil
}

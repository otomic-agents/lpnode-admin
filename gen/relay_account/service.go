// Code generated by goa v3.11.0, DO NOT EDIT.
//
// relayAccount service
//
// Command:
// $ goa gen admin-panel/design

package relayaccount

import (
	"context"
)

// used to manage lp account on relay
type Service interface {
	// ListAccount implements listAccount.
	ListAccount(context.Context) (res *ListAccountResult, err error)
	// RegisterAccount implements registerAccount.
	RegisterAccount(context.Context, *RegisterAccountPayload) (res *RegisterAccountResult, err error)
	// DeleteAccount implements deleteAccount.
	DeleteAccount(context.Context, *DeleteAccountPayload) (res *DeleteAccountResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "relayAccount"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"listAccount", "registerAccount", "deleteAccount"}

// DeleteAccountPayload is the payload type of the relayAccount service
// deleteAccount method.
type DeleteAccountPayload struct {
	ID string
}

// DeleteAccountResult is the result type of the relayAccount service
// deleteAccount method.
type DeleteAccountResult struct {
	Code    *int64
	Result  *string
	Message *string
}

// ListAccountResult is the result type of the relayAccount service listAccount
// method.
type ListAccountResult struct {
	Code    int64
	Result  []*RelayAccountItem
	Message *string
}

// RegisterAccountPayload is the payload type of the relayAccount service
// registerAccount method.
type RegisterAccountPayload struct {
	RelayURL string
	Profile  string
}

// RegisterAccountResult is the result type of the relayAccount service
// registerAccount method.
type RegisterAccountResult struct {
	Code    *int64
	Result  *RelayAccountItem
	Message *string
}

type RelayAccountItem struct {
	ID           *string
	Name         *string
	Profile      *string
	LpIDFake     *string
	LpNodeAPIKey *string
	RelayAPIKey  *string
}

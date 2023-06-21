// Code generated by goa v3.11.0, DO NOT EDIT.
//
// relayAccount HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	relayaccount "admin-panel/gen/relay_account"
	"encoding/json"
	"fmt"
)

// BuildRegisterAccountPayload builds the payload for the relayAccount
// registerAccount endpoint from CLI flags.
func BuildRegisterAccountPayload(relayAccountRegisterAccountBody string) (*relayaccount.RegisterAccountPayload, error) {
	var err error
	var body RegisterAccountRequestBody
	{
		err = json.Unmarshal([]byte(relayAccountRegisterAccountBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"name\": \"Modi aut rerum repellendus qui vel.\",\n      \"profile\": \"Ullam incidunt labore rem quibusdam.\"\n   }'")
		}
	}
	v := &relayaccount.RegisterAccountPayload{
		Name:    body.Name,
		Profile: body.Profile,
	}

	return v, nil
}

// BuildDeleteAccountPayload builds the payload for the relayAccount
// deleteAccount endpoint from CLI flags.
func BuildDeleteAccountPayload(relayAccountDeleteAccountBody string) (*relayaccount.DeleteAccountPayload, error) {
	var err error
	var body DeleteAccountRequestBody
	{
		err = json.Unmarshal([]byte(relayAccountDeleteAccountBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"id\": \"Debitis fugiat dolores asperiores velit illo harum.\"\n   }'")
		}
	}
	v := &relayaccount.DeleteAccountPayload{
		ID: body.ID,
	}

	return v, nil
}

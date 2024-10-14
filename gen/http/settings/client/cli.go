// Code generated by goa v3.11.0, DO NOT EDIT.
//
// settings HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	settings "admin-panel/gen/settings"
	"encoding/json"
	"fmt"
)

// BuildSettingsPayload builds the payload for the settings settings endpoint
// from CLI flags.
func BuildSettingsPayload(settingsSettingsBody string) (*settings.SettingsPayload, error) {
	var err error
	var body SettingsRequestBody
	{
		err = json.Unmarshal([]byte(settingsSettingsBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"relayUri\": \"A est sint autem dolorem voluptas.\"\n   }'")
		}
	}
	v := &settings.SettingsPayload{
		RelayURI: body.RelayURI,
	}

	return v, nil
}

// Code generated by goa v3.11.0, DO NOT EDIT.
//
// configResource HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	configresource "admin-panel/gen/config_resource"
	"encoding/json"
	"fmt"
)

// BuildCreateResourcePayload builds the payload for the configResource
// createResource endpoint from CLI flags.
func BuildCreateResourcePayload(configResourceCreateResourceBody string) (*configresource.CreateResourcePayload, error) {
	var err error
	var body CreateResourceRequestBody
	{
		err = json.Unmarshal([]byte(configResourceCreateResourceBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"appName\": \"Qui vel et qui.\",\n      \"clientId\": \"Et dolor laboriosam unde beatae.\",\n      \"template\": \"Exercitationem aut.\",\n      \"version\": \"Incidunt sit iure.\"\n   }'")
		}
	}
	v := &configresource.CreateResourcePayload{
		AppName:  body.AppName,
		Version:  body.Version,
		ClientID: body.ClientID,
		Template: body.Template,
	}

	return v, nil
}

// BuildGetResourcePayload builds the payload for the configResource
// getResource endpoint from CLI flags.
func BuildGetResourcePayload(configResourceGetResourceBody string) (*configresource.GetResourcePayload, error) {
	var err error
	var body GetResourceRequestBody
	{
		err = json.Unmarshal([]byte(configResourceGetResourceBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"clientId\": \"Blanditiis adipisci cupiditate officiis.\"\n   }'")
		}
	}
	v := &configresource.GetResourcePayload{
		ClientID: body.ClientID,
	}

	return v, nil
}

// BuildEditResultPayload builds the payload for the configResource editResult
// endpoint from CLI flags.
func BuildEditResultPayload(configResourceEditResultBody string) (*configresource.EditResultPayload, error) {
	var err error
	var body EditResultRequestBody
	{
		err = json.Unmarshal([]byte(configResourceEditResultBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"appName\": \"Sed odio ut ratione voluptates est.\",\n      \"clientId\": \"Soluta sapiente.\",\n      \"template\": \"Et sequi et quaerat voluptatem.\",\n      \"templateResult\": \"Officia fugit voluptatem iusto dolore aut omnis.\",\n      \"version\": \"Vitae odio officia neque possimus.\",\n      \"versionHash\": \"Voluptatem consequuntur excepturi omnis voluptatem.\"\n   }'")
		}
	}
	v := &configresource.EditResultPayload{
		TemplateResult: body.TemplateResult,
		Template:       body.Template,
		ClientID:       body.ClientID,
		AppName:        body.AppName,
		Version:        body.Version,
		VersionHash:    body.VersionHash,
	}

	return v, nil
}

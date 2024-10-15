// Code generated by goa v3.11.0, DO NOT EDIT.
//
// installCtrlPanel HTTP client CLI support package
//
// Command:
// $ goa gen admin-panel/design

package client

import (
	installctrlpanel "admin-panel/gen/install_ctrl_panel"
	"encoding/json"
	"fmt"

	goa "goa.design/goa/v3/pkg"
)

// BuildListInstallPayload builds the payload for the installCtrlPanel
// listInstall endpoint from CLI flags.
func BuildListInstallPayload(installCtrlPanelListInstallBody string) (*installctrlpanel.ListInstallPayload, error) {
	var err error
	var body ListInstallRequestBody
	{
		err = json.Unmarshal([]byte(installCtrlPanelListInstallBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"installType\": \"Omnis modi.\"\n   }'")
		}
	}
	v := &installctrlpanel.ListInstallPayload{
		InstallType: body.InstallType,
	}

	return v, nil
}

// BuildInstallLpClientPayload builds the payload for the installCtrlPanel
// installLpClient endpoint from CLI flags.
func BuildInstallLpClientPayload(installCtrlPanelInstallLpClientBody string) (*installctrlpanel.InstallLpClientPayload, error) {
	var err error
	var body InstallLpClientRequestBody
	{
		err = json.Unmarshal([]byte(installCtrlPanelInstallLpClientBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"setupConfig\": {\n         \"awsAccessKeyId\": \"Aliquam vel dolorem deleniti magnam et.\",\n         \"awsSecretAccessKey\": \"Neque voluptatem nihil debitis et magnam nisi.\",\n         \"connectionExplorerurl\": \"Quia sunt quam.\",\n         \"connectionHelperurl\": \"Laboriosam tenetur numquam.\",\n         \"connectionNodeurl\": \"Perferendis ut unde voluptatibus.\",\n         \"connectionWalleturl\": \"Veniam sint neque quia esse fugit.\",\n         \"containerPort\": \"Sint dolorem.\",\n         \"customEnv\": [\n            {\n               \"key\": \"Ipsa et iusto ab.\",\n               \"value\": \"Dignissimos cupiditate.\"\n            },\n            {\n               \"key\": \"Ipsa et iusto ab.\",\n               \"value\": \"Dignissimos cupiditate.\"\n            }\n         ],\n         \"deploymentName\": \"Sint asperiores error nulla quo sunt.\",\n         \"imageRepository\": \"Est magnam sed.\",\n         \"install\": true,\n         \"rpcUrl\": \"Laboriosam laboriosam error.\",\n         \"serviceName\": \"Qui doloremque nam.\",\n         \"startBlock\": \"Sunt dolor ea ducimus doloribus.\",\n         \"type\": \"Sit sed.\"\n      }\n   }'")
		}
		if body.SetupConfig == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("setupConfig", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &installctrlpanel.InstallLpClientPayload{}
	if body.SetupConfig != nil {
		v.SetupConfig = marshalAmmClientSetupConfigRequestBodyToInstallctrlpanelAmmClientSetupConfig(body.SetupConfig)
	}

	return v, nil
}

// BuildUninstallLpClientPayload builds the payload for the installCtrlPanel
// uninstallLpClient endpoint from CLI flags.
func BuildUninstallLpClientPayload(installCtrlPanelUninstallLpClientBody string) (*installctrlpanel.UninstallLpClientPayload, error) {
	var err error
	var body UninstallLpClientRequestBody
	{
		err = json.Unmarshal([]byte(installCtrlPanelUninstallLpClientBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"setupConfig\": {\n         \"type\": \"Nulla magnam.\",\n         \"uninstall\": false\n      }\n   }'")
		}
		if body.SetupConfig == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("setupConfig", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &installctrlpanel.UninstallLpClientPayload{}
	if body.SetupConfig != nil {
		v.SetupConfig = marshalAmmClientUnSetupConfigRequestBodyToInstallctrlpanelAmmClientUnSetupConfig(body.SetupConfig)
	}

	return v, nil
}

// BuildInstallDeploymentPayload builds the payload for the installCtrlPanel
// installDeployment endpoint from CLI flags.
func BuildInstallDeploymentPayload(installCtrlPanelInstallDeploymentBody string) (*installctrlpanel.InstallDeploymentPayload, error) {
	var err error
	var body InstallDeploymentRequestBody
	{
		err = json.Unmarshal([]byte(installCtrlPanelInstallDeploymentBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"setupConfig\": {\n         \"containerPort\": \"Incidunt eos est ipsa aut ratione eum.\",\n         \"customEnv\": [\n            {\n               \"key\": \"Ipsa et iusto ab.\",\n               \"value\": \"Dignissimos cupiditate.\"\n            },\n            {\n               \"key\": \"Ipsa et iusto ab.\",\n               \"value\": \"Dignissimos cupiditate.\"\n            },\n            {\n               \"key\": \"Ipsa et iusto ab.\",\n               \"value\": \"Dignissimos cupiditate.\"\n            },\n            {\n               \"key\": \"Ipsa et iusto ab.\",\n               \"value\": \"Dignissimos cupiditate.\"\n            }\n         ],\n         \"imageRepository\": \"Aliquid quibusdam deserunt aut.\",\n         \"install\": false,\n         \"installType\": \"market\",\n         \"name\": \"Eveniet adipisci quibusdam nihil fugit.\"\n      }\n   }'")
		}
		if body.SetupConfig == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("setupConfig", "body"))
		}
		if body.SetupConfig != nil {
			if err2 := ValidateDeploymentSetupConfigRequestBody(body.SetupConfig); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	v := &installctrlpanel.InstallDeploymentPayload{}
	if body.SetupConfig != nil {
		v.SetupConfig = marshalDeploymentSetupConfigRequestBodyToInstallctrlpanelDeploymentSetupConfig(body.SetupConfig)
	}

	return v, nil
}

// BuildUninstallDeploymentPayload builds the payload for the installCtrlPanel
// uninstallDeployment endpoint from CLI flags.
func BuildUninstallDeploymentPayload(installCtrlPanelUninstallDeploymentBody string) (*installctrlpanel.UninstallDeploymentPayload, error) {
	var err error
	var body UninstallDeploymentRequestBody
	{
		err = json.Unmarshal([]byte(installCtrlPanelUninstallDeploymentBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"setupConfig\": {\n         \"installType\": \"Totam voluptatem.\",\n         \"name\": \"Ex sunt quidem dolores est.\",\n         \"uninstall\": true\n      }\n   }'")
		}
		if body.SetupConfig == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("setupConfig", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &installctrlpanel.UninstallDeploymentPayload{}
	if body.SetupConfig != nil {
		v.SetupConfig = marshalUnDeploymentSetupConfigRequestBodyToInstallctrlpanelUnDeploymentSetupConfig(body.SetupConfig)
	}

	return v, nil
}

// BuildUpdateDeploymentPayload builds the payload for the installCtrlPanel
// updateDeployment endpoint from CLI flags.
func BuildUpdateDeploymentPayload(installCtrlPanelUpdateDeploymentBody string) (*installctrlpanel.UpdateDeploymentPayload, error) {
	var err error
	var body UpdateDeploymentRequestBody
	{
		err = json.Unmarshal([]byte(installCtrlPanelUpdateDeploymentBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"setupConfig\": {\n         \"installContext\": \"Architecto eligendi necessitatibus nisi consequatur illum.\",\n         \"installType\": \"Ut quis.\",\n         \"name\": \"Iure et blanditiis unde beatae soluta.\",\n         \"update\": false\n      }\n   }'")
		}
		if body.SetupConfig == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("setupConfig", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &installctrlpanel.UpdateDeploymentPayload{}
	if body.SetupConfig != nil {
		v.SetupConfig = marshalUpdateDeploymentConfigRequestBodyToInstallctrlpanelUpdateDeploymentConfig(body.SetupConfig)
	}

	return v, nil
}

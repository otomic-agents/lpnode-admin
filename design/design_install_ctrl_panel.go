package design

import (
	. "goa.design/goa/v3/dsl"
)

var Ctrl_DeploayItem = Type("ctrlDeploayItem", func() {
	Attribute("installType", String, "install type")
	Attribute("name", String, "name")
	Attribute("status", Int64, "install status")
	Attribute("installContext", String, "install context")
	Attribute("yaml", String, "yaml")
})
var Ctrl_Amm_Client_SetupConfig = Type("ammClientSetupConfig", func() {
	Attribute("customEnv", ArrayOf(Deployment_SetupConfig_Env_Item))
	Attribute("imageRepository", String, "imageRepository")
	Attribute("serviceName", String, "serviceName")
	Attribute("deploymentName", String, "deploymentName")
	Attribute("type", String, "type") // avax bsc near xrp
	Attribute("startBlock", String, "")
	Attribute("rpcUrl", String, "")
	Attribute("connectionNodeurl", String, "")
	Attribute("connectionWalleturl", String, "")
	Attribute("connectionHelperurl", String, "")
	Attribute("connectionExplorerurl", String, "")
	Attribute("awsAccessKeyId", String, "")
	Attribute("containerPort", String)
	Attribute("awsSecretAccessKey", String, "")
	Attribute("install", Boolean, "")
	Required("imageRepository", "type", "install")
})
var Ctrl_Amm_Client_UnSetupConfig = Type("ammClientUnSetupConfig", func() {
	Attribute("type", String, "") // avax bsc near xrp
	Attribute("uninstall", Boolean, "")
})
var Deployment_SetupConfig_Env_Item = Type("DeploymentSetupConfigEnvItem", func() {
	Attribute("key", String)
	Attribute("value", String)
})
var Ctrl_Deployment_SetupConfig = Type("DeploymentSetupConfig", func() {
	Attribute("imageRepository", String, "")
	Attribute("containerPort", String, "")
	Attribute("install", Boolean, "")
	Attribute("installType", String, func() {
		Enum("ammClient", "market", "amm", "userApp")
	})
	Attribute("name", String, "")
	Attribute("customEnv", ArrayOf(Deployment_SetupConfig_Env_Item), "env list")
	Required("installType", "install", "imageRepository", "name")
})
var Ctrl_UnDeployment_SetupConfig = Type("UnDeploymentSetupConfig", func() {
	Attribute("uninstall", Boolean, "")
	Attribute("installType", String)
	Attribute("name", String)
	Required("installType", "uninstall", "name")
})
var Ctrl_UpdateDeploymentConfig = Type("updateDeploymentConfig", func() {
	Attribute("installType", String)
	Attribute("name", String)
	Attribute("installContext", String)
	Attribute("update", Boolean)
	Required("name", "installType", "update")
})
var Ctrl_InstallDeploymentDataResult = Type("installDeploymentDataResult", func() {
	Attribute("CmdStdout", String)
	Attribute("CmdStderr", String)
	Attribute("Template", String)
})
var Ctrl_UnInstallDeploymentDataResult = Type("unInstallDeploymentDataResult", func() {
	Attribute("CmdStdout", String)
	Attribute("CmdStderr", String)
	Attribute("Template", String)
})
var _ = Service("installCtrlPanel", func() {
	Description("used to control install and startup of nodes")
	Method("listInstall", func() {
		Payload(func() {
			Attribute("installType", String, "type of service installed")
			Required("installType")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(Ctrl_DeploayItem), "list of installed services")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ctrl/listInstall")
		})
	})
	Method("installLpClient", func() {
		Error("an_error")
		Payload(func() {
			Attribute("setupConfig", Ctrl_Amm_Client_SetupConfig, "ammClientInstallConfig")
			Required("setupConfig")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", func() {
				Attribute("Template", String, "rendered template content")
				Attribute("CmdStdout", String, "")
				Attribute("CmdStderr", String, "")
			})
			Attribute("message", String)

		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ctrl/installLpClient")
			Response(StatusOK)
		})
	})
	Method("uninstallLpClient", func() {
		Payload(func() {
			Attribute("setupConfig", Ctrl_Amm_Client_UnSetupConfig, "ammClientUninstallConfig")
			Required("setupConfig")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", func() {
				Attribute("Template", String, "rendered template content")
				Attribute("CmdStdout", String, "")
				Attribute("CmdStderr", String, "")
			})
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ctrl/uninstallLpClient")
		})
	})

	Method("installDeployment", func() {
		Payload(func() {
			Attribute("setupConfig", Ctrl_Deployment_SetupConfig, "deploymentSetupConfig")
			Required("setupConfig")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Ctrl_InstallDeploymentDataResult, "install result")
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ctrl/installDeployment")
		})
	})
	Method("uninstallDeployment", func() {
		Payload(func() {
			Attribute("setupConfig", Ctrl_UnDeployment_SetupConfig, "UnDeploymentSetupConfig")
			Required("setupConfig")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", Ctrl_UnInstallDeploymentDataResult)
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ctrl/uninstallDeployment")
		})
	})
	Method("updateDeployment", func() {
		Payload(func() {
			Attribute("setupConfig", Ctrl_UpdateDeploymentConfig, "updateDeploymentConfig")
			Required("setupConfig")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", func() {
				Attribute("CmdStdout", String, "")
				Attribute("CmdStderr", String, "")
				Attribute("Template", String, "rendered template content")
			})
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ctrl/updateDeployment")
		})
	})
})

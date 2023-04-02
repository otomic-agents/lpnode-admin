package design

import (
	. "goa.design/goa/v3/dsl"
)

var Ctrl_DeploayItem = Type("ctrlDeploayItem", func() {
	Attribute("installType", String, "安装的类型 Client Amm Market等等...")
	Attribute("name", String, "具体的服务名称，如bsc avax ")
	Attribute("status", Int64, "服务当前的安装状态")
	Attribute("installContext", String, "之前安装的上下文")
	Attribute("yaml", String, "安装模版的原始内容")
})
var Ctrl_Amm_Client_SetupConfig = Type("ammClientSetupConfig", func() {
	Attribute("customEnv", ArrayOf(Deployment_SetupConfig_Env_Item))
	Attribute("imageRepository", String, "需要安装的镜像地址")
	Attribute("serviceName", String, "需要安装的ServiceName")
	Attribute("deploymentName", String, "需要安装的deploymentName")
	Attribute("type", String, "安装的类型") // avax bsc near xrp
	Attribute("startBlock", String, "")
	Attribute("rpcUrl", String, "")
	Attribute("connectionNodeurl", String, "")
	Attribute("connectionWalleturl", String, "")
	Attribute("connectionHelperurl", String, "")
	Attribute("connectionExplorerurl", String, "")
	Attribute("awsAccessKeyId", String, "")
	Attribute("containerPort", String)
	Attribute("awsSecretAccessKey", String, "")
	Attribute("install", Boolean, "是否直接安装")
	Required("imageRepository", "type", "install")
})
var Ctrl_Amm_Client_UnSetupConfig = Type("ammClientUnSetupConfig", func() {
	Attribute("type", String, "安装的类型") // avax bsc near xrp
	Attribute("uninstall", Boolean, "是否直接操作")
})
var Deployment_SetupConfig_Env_Item = Type("DeploymentSetupConfigEnvItem", func() {
	Attribute("key", String)
	Attribute("value", String)
})
var Ctrl_Deployment_SetupConfig = Type("DeploymentSetupConfig", func() {
	Attribute("imageRepository", String, "需要安装的镜像地址")
	Attribute("containerPort", String, "容器的端口 可选")
	Attribute("install", Boolean, "是否直接安装")
	Attribute("installType", String, func() {
		Enum("ammClient", "market", "amm", "userApp")
	})
	Attribute("name", String, "这个服务叫什么名字")
	Attribute("customEnv", ArrayOf(Deployment_SetupConfig_Env_Item), "Env配置列表")
	Required("installType", "install", "imageRepository", "name")
})
var Ctrl_UnDeployment_SetupConfig = Type("UnDeploymentSetupConfig", func() {
	Attribute("uninstall", Boolean, "是否直接卸载")
	Attribute("installType", String)
	Attribute("name", String)
	Required("installType", "uninstall", "name")
})
var Ctrl_UpdateDeploymentConfig = Type("updateDeploymentConfig", func() {
	Attribute("installType", String)
	Attribute("name", String)
	Attribute("installContext", String)
	Attribute("update", Boolean) // 是否直接更新
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
	Description("用于控制各个节点的安装和启动")
	Method("listInstall", func() {
		Description("列出已经安装的服务和状态,每一项针对多个服务，针对一个配置文件")
		Payload(func() {
			Attribute("installType", String, "安装的服务类型")
			Required("installType")
		})
		Result(func() {
			Attribute("code", Int64, "")
			Attribute("result", ArrayOf(Ctrl_DeploayItem), "已经安装的列表")
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
				Attribute("Template", String, "渲染后的模版内容")
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
				Attribute("Template", String, "渲染后的模版内容")
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
			Attribute("result", Ctrl_InstallDeploymentDataResult, "安装完成的结果")
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
				Attribute("Template", String, "渲染后的模版内容")
			})
			Attribute("message", String)
		})
		HTTP(func() {
			POST("/lpnode/lpnode_admin_panel/ctrl/updateDeployment")
		})
	})
})

// Code generated by goa v3.11.0, DO NOT EDIT.
//
// installCtrlPanel service
//
// Command:
// $ goa gen admin-panel/design

package installctrlpanel

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// 用于控制各个节点的安装和启动
type Service interface {
	// 列出已经安装的服务和状态,每一项针对多个服务，针对一个配置文件
	ListInstall(context.Context, *ListInstallPayload) (res *ListInstallResult, err error)
	// InstallLpClient implements installLpClient.
	InstallLpClient(context.Context, *InstallLpClientPayload) (res *InstallLpClientResult, err error)
	// UninstallLpClient implements uninstallLpClient.
	UninstallLpClient(context.Context, *UninstallLpClientPayload) (res *UninstallLpClientResult, err error)
	// InstallDeployment implements installDeployment.
	InstallDeployment(context.Context, *InstallDeploymentPayload) (res *InstallDeploymentResult, err error)
	// UninstallDeployment implements uninstallDeployment.
	UninstallDeployment(context.Context, *UninstallDeploymentPayload) (res *UninstallDeploymentResult, err error)
	// UpdateDeployment implements updateDeployment.
	UpdateDeployment(context.Context, *UpdateDeploymentPayload) (res *UpdateDeploymentResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "installCtrlPanel"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [6]string{"listInstall", "installLpClient", "uninstallLpClient", "installDeployment", "uninstallDeployment", "updateDeployment"}

type AmmClientSetupConfig struct {
	CustomEnv []*DeploymentSetupConfigEnvItem
	// 需要安装的镜像地址
	ImageRepository string
	// 需要安装的ServiceName
	ServiceName *string
	// 需要安装的deploymentName
	DeploymentName *string
	// 安装的类型
	Type                  string
	StartBlock            *string
	RPCURL                *string
	ConnectionNodeurl     *string
	ConnectionWalleturl   *string
	ConnectionHelperurl   *string
	ConnectionExplorerurl *string
	AwsAccessKeyID        *string
	ContainerPort         *string
	AwsSecretAccessKey    *string
	// 是否直接安装
	Install bool
}

type AmmClientUnSetupConfig struct {
	// 安装的类型
	Type *string
	// 是否直接操作
	Uninstall *bool
}

type CtrlDeploayItem struct {
	// 安装的类型 Client Amm Market等等...
	InstallType *string
	// 具体的服务名称，如bsc avax
	Name *string
	// 服务当前的安装状态
	Status *int64
	// 之前安装的上下文
	InstallContext *string
	// 安装模版的原始内容
	Yaml *string
}

type DeploymentSetupConfig struct {
	// 需要安装的镜像地址
	ImageRepository string
	// 容器的端口 可选
	ContainerPort *string
	// 是否直接安装
	Install     bool
	InstallType string
	// 这个服务叫什么名字
	Name string
	// Env配置列表
	CustomEnv []*DeploymentSetupConfigEnvItem
}

type DeploymentSetupConfigEnvItem struct {
	Key   *string
	Value *string
}

type InstallDeploymentDataResult struct {
	CmdStdout *string
	CmdStderr *string
	Template  *string
}

// InstallDeploymentPayload is the payload type of the installCtrlPanel service
// installDeployment method.
type InstallDeploymentPayload struct {
	// deploymentSetupConfig
	SetupConfig *DeploymentSetupConfig
}

// InstallDeploymentResult is the result type of the installCtrlPanel service
// installDeployment method.
type InstallDeploymentResult struct {
	Code *int64
	// 安装完成的结果
	Result  *InstallDeploymentDataResult
	Message *string
}

// InstallLpClientPayload is the payload type of the installCtrlPanel service
// installLpClient method.
type InstallLpClientPayload struct {
	// ammClientInstallConfig
	SetupConfig *AmmClientSetupConfig
}

// InstallLpClientResult is the result type of the installCtrlPanel service
// installLpClient method.
type InstallLpClientResult struct {
	Code   *int64
	Result *struct {
		// 渲染后的模版内容
		Template  *string
		CmdStdout *string
		CmdStderr *string
	}
	Message *string
}

// ListInstallPayload is the payload type of the installCtrlPanel service
// listInstall method.
type ListInstallPayload struct {
	// 安装的服务类型
	InstallType string
}

// ListInstallResult is the result type of the installCtrlPanel service
// listInstall method.
type ListInstallResult struct {
	Code *int64
	// 已经安装的列表
	Result  []*CtrlDeploayItem
	Message *string
}

type UnDeploymentSetupConfig struct {
	// 是否直接卸载
	Uninstall   bool
	InstallType string
	Name        string
}

type UnInstallDeploymentDataResult struct {
	CmdStdout *string
	CmdStderr *string
	Template  *string
}

// UninstallDeploymentPayload is the payload type of the installCtrlPanel
// service uninstallDeployment method.
type UninstallDeploymentPayload struct {
	// UnDeploymentSetupConfig
	SetupConfig *UnDeploymentSetupConfig
}

// UninstallDeploymentResult is the result type of the installCtrlPanel service
// uninstallDeployment method.
type UninstallDeploymentResult struct {
	Code    *int64
	Result  *UnInstallDeploymentDataResult
	Message *string
}

// UninstallLpClientPayload is the payload type of the installCtrlPanel service
// uninstallLpClient method.
type UninstallLpClientPayload struct {
	// ammClientUninstallConfig
	SetupConfig *AmmClientUnSetupConfig
}

// UninstallLpClientResult is the result type of the installCtrlPanel service
// uninstallLpClient method.
type UninstallLpClientResult struct {
	Code   *int64
	Result *struct {
		// 渲染后的模版内容
		Template  *string
		CmdStdout *string
		CmdStderr *string
	}
	Message *string
}

type UpdateDeploymentConfig struct {
	InstallType    string
	Name           string
	InstallContext *string
	Update         bool
}

// UpdateDeploymentPayload is the payload type of the installCtrlPanel service
// updateDeployment method.
type UpdateDeploymentPayload struct {
	// updateDeploymentConfig
	SetupConfig *UpdateDeploymentConfig
}

// UpdateDeploymentResult is the result type of the installCtrlPanel service
// updateDeployment method.
type UpdateDeploymentResult struct {
	Code   *int64
	Result *struct {
		CmdStdout *string
		CmdStderr *string
		// 渲染后的模版内容
		Template *string
	}
	Message *string
}

// MakeAnError builds a goa.ServiceError from an error.
func MakeAnError(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "an_error", false, false, false)
}

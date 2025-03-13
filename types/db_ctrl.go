package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type TemplateWriter struct {
	ByteBuffer []byte
}

func (templateWriter *TemplateWriter) Write(p []byte) (n int, err error) {
	for _, v := range p {
		templateWriter.ByteBuffer = append(templateWriter.ByteBuffer, v)
	}
	return len(p), nil
}

type LatestConfig struct {
	Service    ClientSetupConfigService    `bson:"service" json:"service"`
	Deployment ClientSetupConfigDeployment `bson:"deployment" json:"deployment"`
}
type UpdateResult struct {
	LatestConfig LatestConfig `bson:"latestConfig" json:"latestConfig"`
}
type InstallRow struct {
	ID                   primitive.ObjectID `bson:"_id"`
	InstallType          string             `bson:"installType"` // client amm market
	Name                 string             `bson:"name"`        // bsc avax
	Status               int64              `bson:"status"`
	Yaml                 string             `bson:"yaml"`
	Stdout               string             `bson:"stdout"`
	Stderr               string             `bson:"stderr"`
	InstallContext       string             `bson:"installContext"`
	ConfigStatus         int64              `bson:"configStatus"`
	ServiceName          string             `bson:"serviceName"`
	ChainId              int64              `bson:"chainId"`
	UninstallStdErr      string             `bson:"un_stderr"`
	UninstallStdOut      string             `bson:"un_stdout"`
	ChainType            string             `bson:"chainType"`
	Namespace            string             `bson:"namespace"`
	RegisterClientStatus int64              `bson:"registerClientStatus"`
	EnvList              []struct {
		Name  string `bson:"name"`
		Value string `bson:"value"`
	} `bson:"envList"`
	UpdateResult UpdateResult `bson:"updateResult" json:"updateResult"`
}

// Client Setup
type SetupConfig struct {
	Service    ClientSetupConfigService    `json:"service" bson:service`
	Deployment ClientSetupConfigDeployment `json:"deployment" bson:"deployment"`
}
type ClientSetupConfigService struct {
}
type ClientSetupConfigDeployment struct {
	RunEnv                string                              `json:"runEnv" bson:"runEnv"`
	Name                  string                              `json:"name" bson:"name"`
	Namespace             string                              `json:"namespace" bson:"namespace"`
	Image                 string                              `json:"image" bson:"image"`
	StartBlock            string                              `json:"startBlock" bson:"startBlock"`
	RpcUrl                string                              `json:"rpcUrl" bson:"rpcUrl"`
	ConnectionNodeurl     string                              `json:"connectionNodeurl" bson:"connectionNodeurl"`
	ConnectionWalleturl   string                              `json:"connectionWalleturl" bson:"connectionWalleturl"`
	ConnectionHelperurl   string                              `json:"connectionHelperurl" bson:"connectionHelperurl"`
	ConnectionExplorerurl string                              `json:"connectionExplorerurl" bson:"connectionExplorerurl"`
	AwsAccessKeyId        string                              `json:"awsAccessKeyId" bson:"awsAccessKeyId"`
	AwsSecretAccessKey    string                              `json:"awsSecretAccessKey" bson:"awsSecretAccessKey"`
	CustomEnv             []AmmSetupConfigDeploymentCustomEnv `json:"customEnv" bson:"customEnv"`
	ContainerPort         string                              `json:"containerPort" bson:"containerPort"`
	OsSystemServer        string                              `json:"osSystemServer" bson:"osSystemServer"`
	OsApiSecret           string                              `json:"osApiSecret" bson:"osApiSecret"`
	OsApiKey              string                              `json:"osApiKey" bson:"osApiKey"`
	RedisHost             string                              `json:"redisHost" bson:"redisHost"`
	MongodbHost           string                              `json:"mongodbHost" bson:"mongodbHost"`
	MongodbPass           string                              `json:"mongodbPass" bson:"mongodbPass"`
	RedisPass             string                              `json:"redisPass" bson:"redisPass"`
	RedisPort             string                              `json:"redisPort" bson:"redisPort"`
	MongodbPort           string                              `json:"mongodbPort" bson:"mongodbPort"`
	MongodbAccount        string                              `json:"mongodbAccount" bson:"mongodbAccount"`
	MongodbDbnameLpStore  string                              `json:"mongodbDbnameLpStore" bson:"mongodbDbnameLpStore"`
	MongodbDbnameHistory  string                              `json:"mongodbDbnameHistory" bson:"mongodbDbnameHistory"`
}
type AmmSetupConfigDeploymentCustomEnv struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

// amm setup
type AmmSetupConfigService struct {
}
type AmmSetupConfigDeployment struct {
	Namespace            string                              `json:"namespace"`
	CustomEnv            []AmmSetupConfigDeploymentCustomEnv `json:"customEnv"`
	Image                string                              `json:"image"`
	Name                 string                              `json:"name"`
	ContainerPort        string                              `json:"containerPort"`
	OsSystemServer       string                              `json:"osSystemServer"`
	OsApiSecret          string                              `json:"osApiSecret"`
	OsApiKey             string                              `json:"osApiKey"`
	RedisHost            string                              `json:"redisHost"`
	MongodbHost          string                              `json:"mongodbHost"`
	MongodbPass          string                              `json:"mongodbPass"`
	RedisPass            string                              `json:"redisPass"`
	RedisPort            string                              `json:"redisPort"`
	MongodbPort          string                              `json:"mongodbPort"`
	MongodbAccount       string                              `json:"mongodbAccount"`
	MongodbDbnameLpStore string                              `json:"mongodbDbnameLpStore"`
	MongodbDbnameHistory string                              `json:"mongodbDbnameHistory"`
}

type AmmSetupConfig struct {
	Service    AmmSetupConfigService    `json:"service"`
	Deployment AmmSetupConfigDeployment `json:"deployment"`
}

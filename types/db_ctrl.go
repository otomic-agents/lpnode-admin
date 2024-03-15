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

type InstallRow struct {
	ID                   primitive.ObjectID `bson:"_id"`
	InstallType          string             `bson:"installType"` // client  amm market
	Name                 string             `bson:"name"`        // bsc avax
	Status               int64              `bson:"status"`
	Yaml                 string             `bson:"yaml"`
	Stdout               string             `bson:"stdout"`
	Stderr               string             `bson:"stderr"`
	InstallContext       string             `bson:"installContext"`
	ConfigStatus         int64              `bson:"configStatus"` //0
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
}

// Client Setup
type SetupConfig struct {
	Service    ClientSetupConfigService    `json:"service"`
	Deployment ClientSetupConfigDeployment `json:"deployment"`
}
type ClientSetupConfigService struct {
}
type ClientSetupConfigDeployment struct {
	RunEnv                string                              `json:"runEnv"`
	Name                  string                              `json:"name"`
	Namespace             string                              `json:"namespace"`
	Image                 string                              `json:"image"`
	StartBlock            string                              `json:"startBlock"`
	RpcUrl                string                              `json:"rpcUrl"`
	ConnectionNodeurl     string                              `json:"connectionNodeurl"`
	ConnectionWalleturl   string                              `json:"connectionWalleturl"`
	ConnectionHelperurl   string                              `json:"connectionHelperurl"`
	ConnectionExplorerurl string                              `json:"connectionExplorerurl"`
	AwsAccessKeyId        string                              `json:"awsAccessKeyId"`
	AwsSecretAccessKey    string                              `json:"awsSecretAccessKey"`
	CustomEnv             []AmmSetupConfigDeploymentCustomEnv `json:"customEnv"`
	ContainerPort         string                              `json:"containerPort"`
	OsSystemServer        string                              `json:"osSystemServer"`
	OsApiSecret           string                              `json:"osApiSecret"`
	OsApiKey              string                              `json:"osApiKey"`
	RedisHost             string                              `json:"redisHost"`
	MongodbHost           string                              `json:"mongodbHost"`
	MongodbPass           string                              `json:"mongodbPass"`
	RedisPass             string                              `json:"redisPass"`
	RedisPort             string                              `json:"redisPort"`
	MongodbPort           string                              `json:"mongodbPort"`
	MongodbAccount        string                              `json:"mongodbAccount"`
	MongodbDbnameLpStore  string                              `json:"mongodbDbnameLpStore"`
	MongodbDbnameHistory  string                              `json:"mongodbDbnameHistory"`
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

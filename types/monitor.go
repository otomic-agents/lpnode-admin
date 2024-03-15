package types

type MonitorSetupConfig struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Cron       string `json:"cron"`
	ScriptPath string `json:"scriptPath"`

	MongoDbHost     string `json:"mongoDbHost"`
	MongoDbPort     string `json:"mongoDbPort"`
	MongoDbUser     string `json:"mongoDbUser"`
	MongoDbPass     string `json:"mongoDbPass"`
	MongoDbLpStore  string `json:"mongoDbLpStore"`
	MongoDbBusiness string `json:"mongoDbBusiness"`
	RedisHost       string `json:"redisHost"`
	RedisPort       string `json:"redisPort"`
	RedisPass       string `json:"redisPass"`
	RedisDb         int    `json:"redisDb"`
	UserSpacePath   string `json:"userSpacePath"`
}

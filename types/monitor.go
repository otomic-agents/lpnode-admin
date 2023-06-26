package types

type MonitorSetupConfig struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Cron       string `json:"cron"`
	ScriptPath string `json:"scriptPath"`
}

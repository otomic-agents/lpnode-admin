package types

type MonitorSetupConfig struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Corn       string `json:"corn"`
	ScriptPath string `json:"scriptPath"`
}

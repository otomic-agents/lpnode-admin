package types

type SettingsConfig struct {
	RelayUri string `json:"relayUri" bson:"relayUri" env:"RELAY_URI"`
}

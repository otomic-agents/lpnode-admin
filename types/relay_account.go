package types

type RelayRegisterResponse struct {
	Code         int64  `json:"code"`
	Message      string `json:"message"`
	LpIdFake     string `json:"lp_id_fake"`
	Name         string `json:"name"`
	Profile      string `json:"profile"`
	LpnodeApiKey string `json:"lpnode_api_key"`
	RelayApiKey  string `json:"relay_api_key"`
}

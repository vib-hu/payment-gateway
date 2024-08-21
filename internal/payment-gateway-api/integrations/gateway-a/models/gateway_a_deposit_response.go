package models

type GatewayADepositResponse struct {
	IsSuccessful        bool   `json:"is_successful"`
	ResponseCode        string `json:"response_code"`
	ResponseDescription string `json:"response_description"`
}

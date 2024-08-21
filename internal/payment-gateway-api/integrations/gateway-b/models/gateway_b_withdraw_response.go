package models

import "encoding/xml"

type GatewayBWithdrawEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    GatewayBWithdrawBody
}
type GatewayBWithdrawBody struct {
	XMLName                  xml.Name                 `xml:"Body"`
	GatewayBWithdrawResponse GatewayBWithdrawResponse `xml:"WithdrawResponse"`
}
type GatewayBWithdrawResponse struct {
	IsSuccessful        bool   `xml:"is_successful"`
	ResponseCode        string `xml:"response_code"`
	ResponseDescription string `xml:"response_description"`
}

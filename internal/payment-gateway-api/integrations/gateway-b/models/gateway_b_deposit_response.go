package models

import "encoding/xml"

type GatewayBDepositEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    GatewayBDepositBody
}
type GatewayBDepositBody struct {
	XMLName                 xml.Name                `xml:"Body"`
	GatewayBDepositResponse GatewayBDepositResponse `xml:"DepositResponse"`
}
type GatewayBDepositResponse struct {
	IsSuccessful        bool   `xml:"is_successful"`
	ResponseCode        string `xml:"response_code"`
	ResponseDescription string `xml:"response_description"`
}

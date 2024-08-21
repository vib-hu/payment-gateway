package formatters

import (
	"encoding/xml"
)

type SoapFormatter struct{}

func NewSoapFormatter() *SoapFormatter {
	return &SoapFormatter{}
}

func (sf *SoapFormatter) Marshal(data interface{}) ([]byte, error) {
	envelope := struct {
		XMLName      xml.Name `xml:"soapenv:Envelope"`
		XmlnsSoapenv string   `xml:"xmlns:soapenv,attr"`
		Body         struct {
			XMLName xml.Name    `xml:"soapenv:Body"`
			Content interface{} `xml:",omitempty"`
		}
	}{
		XmlnsSoapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: struct {
			XMLName xml.Name    `xml:"soapenv:Body"`
			Content interface{} `xml:",omitempty"`
		}{
			Content: data,
		},
	}

	return xml.MarshalIndent(envelope, "", "  ")
}

func (sf *SoapFormatter) Unmarshal(data []byte, envelope interface{}) error {
	return xml.Unmarshal(data, &envelope)
}

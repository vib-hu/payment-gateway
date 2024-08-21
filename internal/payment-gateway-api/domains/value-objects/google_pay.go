package valueobjects

import "strings"

type GooglePay struct {
	Token string
}

func (googlePay *GooglePay) IsInvalid() bool {
	return googlePay == nil || strings.TrimSpace(googlePay.Token) == ""
}

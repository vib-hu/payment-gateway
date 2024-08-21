package valueobjects

import "strings"

type ApplePay struct {
	Token string
}

func (applePay *ApplePay) IsInvalid() bool {
	return applePay == nil || strings.TrimSpace(applePay.Token) == ""
}

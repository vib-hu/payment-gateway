package valueobjects

import "strings"

type BankAccount struct {
	AccountNumber     string
	RoutingNumber     string
	AccountHolderName string
}

func (bankAccount *BankAccount) IsInvalid() bool {
	return bankAccount == nil || strings.TrimSpace(bankAccount.AccountNumber) == "" ||
		strings.TrimSpace(bankAccount.RoutingNumber) == "" || strings.TrimSpace(bankAccount.AccountHolderName) == ""
}

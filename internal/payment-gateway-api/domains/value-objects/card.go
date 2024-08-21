package valueobjects

import "strings"

type Card struct {
	CardNumber string
	Cvv        string
}

func (card *Card) IsInvalid() bool {
	return card == nil || strings.TrimSpace(card.CardNumber) == "" || strings.TrimSpace(card.CardNumber) == ""
}

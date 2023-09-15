// Поддержка моделей данных.
package models

import (
	"fmt"
)

// Card структра для хранения банковских карт.
type Card struct {
	Owner           string
	ExpirationMonth uint8
	ExpirationYear  uint16
	Bank            string
	Number          uint64
	CVV             uint16
}

// NewCard returns Card.
func NewCard(owner, bank string, number uint64, expMonth uint8, expYear, cvv uint16) *Card {
	return &Card{
		Owner:           owner,
		Bank:            bank,
		Number:          number,
		ExpirationYear:  expYear,
		ExpirationMonth: expMonth,
		CVV:             cvv,
	}
}

// String prints Card.
func (c *Card) String() string {
	return fmt.Sprintf("%s : %s : %d : %d/%d", c.Owner, c.Bank, c.Number, c.ExpirationMonth, c.ExpirationYear)
}

// Show shows card information.
func (c *Card) Show() string {
	return fmt.Sprintf(" %s\n %d\n valid %d/%d  CVV %d\n %s", c.Bank, c.Number, c.ExpirationMonth, c.ExpirationYear, c.CVV, c.Owner)
}

// IsReadyForStorage check Card's fields for ready to store.
func (c *Card) IsReadyForStorage() (bool, string) {
	const (
		wrongCardOwner           = "wrong owner"
		wrongCardExpirationMonth = "wrong month of expiration (should be from 1 to 12)"
		wrongCardExpirationYear  = "wrong year of expiration (should be more then 2000)"
		wrongCardNumber          = "wrong card number"
		wrongCardCVV             = "wrong card CVV"
	)
	if c.Owner == "" {
		return false, wrongCardOwner
	}
	if c.ExpirationMonth == 0 || c.ExpirationMonth > 12 {
		return false, wrongCardExpirationMonth
	}
	if c.ExpirationYear < 2000 {
		return false, wrongCardExpirationYear
	}
	if c.Number == 0 {
		return false, wrongCardNumber
	}
	if c.CVV == 0 {
		return false, wrongCardCVV
	}
	return true, ""
}

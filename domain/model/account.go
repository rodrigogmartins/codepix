package model

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/asaskevich/govalidator"
)

// Account structure
type Account struct {
	Base      `valid:"required"`
	OwnerName string    `json:"owner_name" valid:"notnull"`
	Bank      *Bank     `valid:"-"`
	Number    string    `json:"number" valid:"notnull"`
	PixKeys   []*PixKey `valid:"-"`
}

func (account *Account) isValid() error {
	_, err := govalidator.ValidateStruct(account)

	if err != nil {
		return err
	}

	return nil
}

// NewAccount instance
func NewAccount(bank *Bank, number string, ownerName string) (*Account, error) {
	account := Account{
		OwnerName: ownerName,
		Bank:      bank,
		Number:    number,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	err := account.isValid()

	if err != nil {
		return nil, err
	}

	return &account, nil
}

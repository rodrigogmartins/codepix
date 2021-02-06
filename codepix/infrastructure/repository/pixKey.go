package repository

import (
	"fmt"

	"github.com/codeedu/imersao/codepix-go/domain/model"
	"github.com/jinzhu/gorm"
)

type PixKeyRepositoryDb struct {
	Db *gorm.DB
}

func (pixKeyRepositoryDb PixKeyRepositoryDb) AddBank(bank *model.Bank) error {
	err := pixKeyRepositoryDb.Db.Create(bank).Error

	if err != nil {
		return err
	}

	return nil
}

func (pixKeyRepositoryDb PixKeyRepositoryDb) FindBank(id string) (*model.Bank, error) {
	var bank model.Bank

	pixKeyRepositoryDb.Db.First(&bank, "id = ?", id)

	if bank.ID == "" {
		return nil, fmt.Errorf("no bank found")
	}

	return &bank, nil
}

func (pixKeyRepositoryDb PixKeyRepositoryDb) AddAccount(account *model.Account) error {
	err := pixKeyRepositoryDb.Db.Create(account).Error

	if err != nil {
		return err
	}

	return nil
}

func (pixKeyRepositoryDb PixKeyRepositoryDb) FindAccount(id string) (*model.Account, error) {
	var account model.Account

	pixKeyRepositoryDb.Db.Preload("Bank").First(&account, "id = ?", id)

	if account.ID == "" {
		return nil, fmt.Errorf("no account found")
	}

	return &account, nil
}

func (pixKeyRepositoryDb PixKeyRepositoryDb) RegisterKey(pixKey *model.PixKey) (*model.PixKey, error) {
	err := pixKeyRepositoryDb.Db.Create(pixKey).Error

	if err != nil {
		return nil, err
	}

	return pixKey, nil
}

func (pixKeyRepositoryDb PixKeyRepositoryDb) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	var pixKey model.PixKey

	pixKeyRepositoryDb.Db.Preload("Account.Bank").First(&pixKey, "kind = ? and key = ?", kind, key)

	if pixKey.ID == "" {
		return nil, fmt.Errorf("no key was found")
	}

	return &pixKey, nil
}

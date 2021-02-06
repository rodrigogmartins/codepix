package repository

import (
	"fmt"

	"github.com/codeedu/imersao/codepix-go/domain/model"
	"github.com/jinzhu/gorm"
)

type TransactionRepositoryDb struct {
	Db *gorm.DB
}

func (transactionRepositoryDb *TransactionRepositoryDb) Register(transaction *model.Transaction) error {
	err := transactionRepositoryDb.Db.Create(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

func (transactionRepositoryDb *TransactionRepositoryDb) Save(transaction *model.Transaction) error {
	err := transactionRepositoryDb.Db.Save(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

func (transactionRepositoryDb *TransactionRepositoryDb) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction

	transactionRepositoryDb.Db.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)

	if transaction.ID == "" {
		return nil, fmt.Errorf("no transaction found")
	}

	return &transaction, nil
}

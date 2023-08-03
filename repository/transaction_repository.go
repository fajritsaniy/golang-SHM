package repository

import (
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	Create(payload *model.Transaction) error
	List() ([]model.Transaction, error)
	Get(id string) (model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func (t *transactionRepository) Create(payload *model.Transaction) error {
	if err := t.db.Omit(clause.Associations).Create(payload).Error; err != nil {
		return err
	}
	return nil
}

func (t *transactionRepository) List() ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := t.db.
		Preload("Vehicle").
		Preload("Customer").
		Preload("Employee").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *transactionRepository) Get(id string) (model.Transaction, error) {
	var transaction model.Transaction
	if err := t.db.Preload(clause.Associations).Where("transaction.id=?", id).First(&transaction).Error; err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

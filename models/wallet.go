package models

import "time"

type WalletModel struct {
	ID        int
	UserID    int
	Balance   float64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionModel struct {
	ID              int
	FromWalletID    int
	ToWalletID      int
	Amount          float64
	TransactionType string
	Status          string
	Description     string
	CreatedAt       time.Time
}

func (w *WalletModel) HasSufficientBalance(amount float64) bool {
	if w.Balance >= amount {
		return true
	}
	return false
}

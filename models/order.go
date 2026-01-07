package models

import "time"

type OrderModel struct {
	ID         int
	UserID     int
	ProductID  int
	Quantity   int
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
}

func (o *OrderModel) IsPending() bool {
	if o.Status == "pending" {
		return true
	}
	return false
}

func (o *OrderModel) IsCompleted() bool {
	if o.Status == "completed" {
		return true
	}
	return false
}

func (o *OrderModel) IsCancelled() bool {
	if o.Status == "cancelled" {
		return true
	}
	return false
}

func (o *OrderModel) CalculateTotal(price float64) float64 {
	return price * float64(o.Quantity)
}

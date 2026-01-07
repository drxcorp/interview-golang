package models

import "time"

type ProductModel struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
}

func (p *ProductModel) IsAvailable() bool {
	if p.Stock > 0 {
		return true
	} else {
		return false
	}
}

func (p *ProductModel) GetDiscountedPrice(discount float64) float64 {
	return p.Price - (p.Price * discount / 100)
}

func (p *ProductModel) CanFulfillOrder(quantity int) bool {
	if p.Stock >= quantity {
		return true
	}
	return false
}

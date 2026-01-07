package models

import "time"

type UserModel struct {
	ID        int
	Name      string
	Email     string
	Password  string
	IsActive  bool
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserModel) Validate() bool {
	if u.Name == "" || u.Email == "" || u.Password == "" {
		return false
	}
	return true
}

func (u *UserModel) IsAdmin() bool {
	if u.Role == "admin" {
		return true
	}
	return false
}

func (u *UserModel) CanCreateOrder() bool {
	return u.IsActive
}

package entity

import (
	"time"

	"github.com/diegoclair/go_boilerplate/util/number"
)

type Account struct {
	ID        int64
	UUID      string
	Name      string
	CPF       string
	Balance   float64
	Password  string
	CreatedAT time.Time
	Active    bool
}

func (a *Account) AddBalance(amount float64) {
	a.Balance = number.RoundFloat(a.Balance+amount, 2)
}

func (a *Account) SubtractBalance(amount float64) {
	a.Balance = number.RoundFloat(a.Balance-amount, 2)
}

func (a *Account) HasSufficientFunds(amount float64) bool {
	return a.Balance >= amount
}

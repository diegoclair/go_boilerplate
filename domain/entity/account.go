package entity

import "time"

type Account struct {
	ID        int64
	UUID      string
	Name      string
	CPF       string `crypt:"true"`
	Balance   float64
	Secret    string
	CreatedAT time.Time
}

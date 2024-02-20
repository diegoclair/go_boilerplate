package account

import "testing"

func TestAccount_AddBalance(t *testing.T) {
	tests := []struct {
		name           string
		initialBalance float64
		amount         float64
		expectedResult float64
	}{
		{
			name:           "Should add balance to account",
			initialBalance: 10.0,
			amount:         5.0,
			expectedResult: 15.0,
		},
		{
			name:           "Should not change balance if amount is zero",
			initialBalance: 10.0,
			amount:         0.0,
			expectedResult: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Balance: tt.initialBalance,
			}

			account.AddBalance(tt.amount)

			if account.Balance != tt.expectedResult {
				t.Errorf("Account.AddBalance() = %v, expected %v", account.Balance, tt.expectedResult)
			}
		})
	}
}

func TestAccount_SubtractBalance(t *testing.T) {
	tests := []struct {
		name           string
		initialBalance float64
		amount         float64
		expectedResult float64
	}{
		{
			name:           "Should subtract balance from account",
			initialBalance: 10.0,
			amount:         5.0,
			expectedResult: 5.0,
		},
		{
			name:           "Should not change balance if amount is zero",
			initialBalance: 10.0,
			amount:         0.0,
			expectedResult: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Balance: tt.initialBalance,
			}

			account.SubtractBalance(tt.amount)

			if account.Balance != tt.expectedResult {
				t.Errorf("Account.SubtractBalance() = %v, expected %v", account.Balance, tt.expectedResult)
			}
		})
	}
}

func TestAccount_HasSufficientFunds(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		amount  float64
		want    bool
	}{
		{
			name: "Should return true when account balance is greater than or equal to the amount",
			account: Account{
				Balance: 100.0,
			},
			amount: 50.0,
			want:   true,
		},
		{
			name: "Should return false when account balance is less than the amount",
			account: Account{
				Balance: 100.0,
			},
			amount: 150.0,
			want:   false,
		},
		{
			name: "Should return true when account balance is equal to the amount",
			account: Account{
				Balance: 100.0,
			},
			amount: 100.0,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.HasSufficientFunds(tt.amount)
			if got != tt.want {
				t.Errorf("Account.HasSufficientFunds() = %v, want %v", got, tt.want)
			}
		})
	}
}

package random

import (
	"testing"
)

func TestRandomName(t *testing.T) {
	tests := []struct {
		name       string
		wantLength int
	}{
		{
			name:       "Should return a random name with 6 characters",
			wantLength: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomName()
			if len(got) != tt.wantLength {
				t.Errorf("RandomName() length = %v, want %v", len(got), tt.wantLength)
			}
		})
	}
}

func TestRandomCPF(t *testing.T) {
	tests := []struct {
		name       string
		wantLength int
	}{
		{
			name:       "Should return a random CPF with 11 characters",
			wantLength: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomCPF()
			if len(got) != tt.wantLength {
				t.Errorf("RandomCPF() length = %v, want %v", len(got), tt.wantLength)
			}
		})
	}
}

func TestRandomSecret(t *testing.T) {
	tests := []struct {
		name       string
		wantLength int
	}{
		{
			name:       "Should return a random secret with 8 characters",
			wantLength: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomPassword()
			if len(got) != tt.wantLength {
				t.Errorf("RandomSecret() length = %v, want %v", len(got), tt.wantLength)
			}
		})
	}
}

func TestRandomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name       string
		args       args
		wantLength int
	}{
		{
			name:       "Should return a random string with 8 characters",
			args:       args{n: 8},
			wantLength: 8,
		},
		{
			name:       "Should return a random string with 16 characters",
			args:       args{n: 16},
			wantLength: 16,
		},
		{
			name:       "Should return a random string with 32 characters",
			args:       args{n: 32},
			wantLength: 32,
		},
		{
			name:       "Should return a random string with 64 characters",
			args:       args{n: 64},
			wantLength: 64,
		},
		{
			name:       "Should return a random string with 128 characters",
			args:       args{n: 128},
			wantLength: 128,
		},
		{
			name:       "Should return a random string with 256 characters",
			args:       args{n: 256},
			wantLength: 256,
		},
		{
			name:       "Should return a random string with 512 characters",
			args:       args{n: 512},
			wantLength: 512,
		},
		{
			name:       "Should return a random string with 1024 characters",
			args:       args{n: 1024},
			wantLength: 1024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomString(tt.args.n)
			if len(got) != tt.wantLength {
				t.Errorf("RandomString() length = %v, want %v", len(got), tt.wantLength)
			}
		})
	}
}

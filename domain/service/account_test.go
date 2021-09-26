package service

import (
	"testing"
)

func Test_accountService_getHashedPassword(t *testing.T) {
	type fields struct {
		svc *Service
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "correct password generated",
			args: args{
				password: "abcd12345678",
			},
			want: "8bf75d25716494a5e1ae63de79db741a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &accountService{
				svc: tt.fields.svc,
			}
			if got := s.getHashedPassword(tt.args.password); got != tt.want {
				t.Errorf("accountService.getHashedPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

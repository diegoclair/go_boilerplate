package errors

import "testing"

func TestSQLNotFound(t *testing.T) {
	type args struct {
		err string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true if the error is because there are no sql rows",
			args: args{err: "no rows in result set"},
			want: true,
		},
		{
			name: "Should return true if the error is because there are no records find with the parameters",
			args: args{err: "No records find"},
			want: true,
		},
		{
			name: "Should return false if the error is not because there are no sql rows or no records find with the parameters",
			args: args{err: "Error"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SQLNotFound(tt.args.err); got != tt.want {
				t.Errorf("SQLNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

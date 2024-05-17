package database

import (
	"testing"

	"github.com/alexver/golang_database/internal/query"
)

func TestExitProcessor_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name EXIT",
			want: "EXIT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewExitProcessor()
			if got := p.Name(); got != tt.want {
				t.Errorf("ExitProcessor.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExitProcessor_Suports(t *testing.T) {
	type args struct {
		query *query.Query
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no support",
			args: args{query: query.CreateQuery("TEST", []string{})},
			want: false,
		},
		{
			name: "ok",
			args: args{query: query.CreateQuery("EXIT", []string{})},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewExitProcessor()
			if got := p.Suports(tt.args.query); got != tt.want {
				t.Errorf("ExitProcessor.Suports() = %v, want %v", got, tt.want)
			}
		})
	}
}

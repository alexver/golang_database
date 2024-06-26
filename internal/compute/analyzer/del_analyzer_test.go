package analyzer

import (
	"testing"

	"github.com/alexver/golang_database/internal/query"
)

func TestDel_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name",
			want: "DEL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel()
			if got := d.Name(); got != tt.want {
				t.Errorf("Del.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDel_Description(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check description",
			want: "Delete value by key from storage. DELETE is an alias. Usage: DEL|DELETE key.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel()
			if got := d.Description(); got != tt.want {
				t.Errorf("Del.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDel_Usage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check usage",
			want: "DEL|DELETE key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel()
			if got := d.Usage(); got != tt.want {
				t.Errorf("Del.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDel_Supports(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "check TEST",
			args: args{name: "TEST"},
			want: false,
		},
		{
			name: "check DEL",
			args: args{name: "DEL"},
			want: true,
		},
		{
			name: "check DELETE",
			args: args{name: "DELETE"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel()
			if got := d.Supports(tt.args.name); got != tt.want {
				t.Errorf("Del.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDel_Validate(t *testing.T) {
	type args struct {
		query *query.Query
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		errString string
	}{
		{
			name:      "fail because wrong command name",
			args:      args{query: query.CreateQuery("ANY", []string{"Test"})},
			wantErr:   true,
			errString: "analyzer DEL error: cannot process 'ANY' command",
		},
		{
			name:      "fail because wrong argument count",
			args:      args{query: query.CreateQuery("DELETE", []string{"Test", "Check"})},
			wantErr:   true,
			errString: "analyzer DEL error: invalid argument count 2",
		},
		{
			name:      "fail because invalid argument",
			args:      args{query: query.CreateQuery("DEL", []string{"Test&&&?"})},
			wantErr:   true,
			errString: "analyzer DEL error: invalid argument #1: Test&&&?",
		},
		{
			name:      "ok",
			args:      args{query: query.CreateQuery("DEL", []string{"Test"})},
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel()
			err := d.Validate(tt.args.query)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errString) {
				t.Errorf("Del.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDel_NormalizeQuery(t *testing.T) {
	type args struct {
		query *query.Query
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check normalization",
			args: args{query: query.CreateQuery("TEST", []string{})},
			want: "DEL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel()
			got := d.NormalizeQuery(tt.args.query)
			if got.GetCommand() != tt.want {
				t.Errorf("Del.NormalizeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

package analyzer

import (
	"testing"

	"github.com/alexver/golang_database/internal/query"
)

func TestPing_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name",
			want: "PING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewPing()
			if got := c.Name(); got != tt.want {
				t.Errorf("Ping.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing_Description(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check description",
			want: "Command to ping test database. Standard answer of the server is PONG.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewPing()
			if got := c.Description(); got != tt.want {
				t.Errorf("Ping.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing_Usage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check usage",
			want: "PING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewPing()
			if got := c.Usage(); got != tt.want {
				t.Errorf("Ping.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing_Supports(t *testing.T) {
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
			name: "check PING",
			args: args{name: "PING"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewPing()
			if got := c.Supports(tt.args.name); got != tt.want {
				t.Errorf("Ping.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing_Validate(t *testing.T) {
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
			errString: "analyzer PING error: cannot process 'ANY' command",
		},
		{
			name:      "fail because wrong argument count",
			args:      args{query: query.CreateQuery("PING", []string{"Test"})},
			wantErr:   true,
			errString: "analyzer PING error: invalid argument count 1",
		},
		{
			name:      "ok",
			args:      args{query: query.CreateQuery("PING", []string{})},
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewPing()
			err := c.Validate(tt.args.query)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errString) {
				t.Errorf("Ping.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPing_NormalizeQuery(t *testing.T) {
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
			want: "PING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewPing()
			got := c.NormalizeQuery(tt.args.query)
			if got.GetCommand() != tt.want {
				t.Errorf("Ping.NormalizeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

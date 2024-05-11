package analyzer

import (
	"testing"

	"github.com/alexver/golang_database/internal/compute/parser"
)

func TestExit_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name",
			want: "EXIT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewExit()
			if got := c.Name(); got != tt.want {
				t.Errorf("Exit.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExit_Description(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check description",
			want: "Command to stop and close test database. You can use QUIT as an alias of EXIT command.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewExit()
			if got := c.Description(); got != tt.want {
				t.Errorf("Exit.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExit_Usage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check usage",
			want: "EXIT|QUIT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewExit()
			if got := c.Usage(); got != tt.want {
				t.Errorf("Exit.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExit_Supports(t *testing.T) {
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
			name: "check EXIT",
			args: args{name: "EXIT"},
			want: true,
		},
		{
			name: "check QUIT",
			args: args{name: "QUIT"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewExit()
			if got := c.Supports(tt.args.name); got != tt.want {
				t.Errorf("Exit.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExit_Validate(t *testing.T) {
	type args struct {
		query parser.Query
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		errString string
	}{
		{
			name:      "fail because wrong command name",
			args:      args{query: parser.CreateQuery("ANY", []string{"Test"})},
			wantErr:   true,
			errString: "analyzer EXIT error: cannot process 'ANY' command",
		},
		{
			name:      "fail because wrong argument count",
			args:      args{query: parser.CreateQuery("QUIT", []string{"Test"})},
			wantErr:   true,
			errString: "analyzer EXIT error: invalid argument count 1",
		},
		{
			name:      "ok",
			args:      args{query: parser.CreateQuery("EXIT", []string{})},
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewExit()
			err := c.Validate(tt.args.query)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errString) {
				t.Errorf("Exit.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

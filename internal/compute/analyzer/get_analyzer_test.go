package analyzer

import (
	"testing"

	"github.com/alexver/golang_database/internal/compute/parser"
)

func TestGet_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name",
			want: "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGet()
			if got := g.Name(); got != tt.want {
				t.Errorf("Get.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet_Description(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check description",
			want: "Get saved value by key. Usage: GET key.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGet()
			if got := g.Description(); got != tt.want {
				t.Errorf("Get.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet_Usage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check usage",
			want: "GET key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGet()
			if got := g.Usage(); got != tt.want {
				t.Errorf("Get.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet_Supports(t *testing.T) {
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
			name: "check GET",
			args: args{name: "GET"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGet()
			if got := g.Supports(tt.args.name); got != tt.want {
				t.Errorf("Get.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet_Validate(t *testing.T) {
	type args struct {
		query *parser.Query
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
			errString: "analyzer GET error: cannot process 'ANY' command",
		},
		{
			name:      "fail because wrong argument count",
			args:      args{query: parser.CreateQuery("GET", []string{"Test", "Check"})},
			wantErr:   true,
			errString: "analyzer GET error: invalid argument count 2",
		},
		{
			name:      "fail because invalid argument",
			args:      args{query: parser.CreateQuery("GET", []string{"Русский"})},
			wantErr:   true,
			errString: "analyzer GET error: invalid argument #1: Русский",
		},
		{
			name:      "ok",
			args:      args{query: parser.CreateQuery("GET", []string{"Test"})},
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGet()
			if err := g.Validate(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("Get.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet_NormalizeQuery(t *testing.T) {
	type args struct {
		query *parser.Query
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check normalization",
			args: args{query: parser.CreateQuery("TEST", []string{})},
			want: "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGet()
			got := g.NormalizeQuery(tt.args.query)
			if got.GetCommand() != tt.want {
				t.Errorf("Get.NormalizeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

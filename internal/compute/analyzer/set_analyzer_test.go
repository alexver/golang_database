package analyzer

import (
	"testing"

	"github.com/alexver/golang_database/internal/query"
)

func TestSet_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name",
			want: "SET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			if got := s.Name(); got != tt.want {
				t.Errorf("Set.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Description(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check description",
			want: "Set value by key. Usage: SET key value.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			if got := s.Description(); got != tt.want {
				t.Errorf("Set.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Usage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check usage",
			want: "SET key value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			if got := s.Usage(); got != tt.want {
				t.Errorf("Set.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Supports(t *testing.T) {
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
			name: "check SET",
			args: args{name: "SET"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			if got := s.Supports(tt.args.name); got != tt.want {
				t.Errorf("Set.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Validate(t *testing.T) {
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
			errString: "analyzer SET error: cannot process 'ANY' command",
		},
		{
			name:      "fail because wrong argument count",
			args:      args{query: query.CreateQuery("SET", []string{"Test"})},
			wantErr:   true,
			errString: "analyzer SET error: invalid argument count 1",
		},
		{
			name:      "fail because invalid argument #1",
			args:      args{query: query.CreateQuery("SET", []string{"===", "value"})},
			wantErr:   true,
			errString: "analyzer SET error: invalid argument #1: ===",
		},
		{
			name:      "fail because invalid argument #2",
			args:      args{query: query.CreateQuery("SET", []string{"test", "wrong!value"})},
			wantErr:   true,
			errString: "analyzer SET error: invalid argument #2: wrong!value",
		},
		{
			name:      "ok",
			args:      args{query: query.CreateQuery("SET", []string{"test", "check_value"})},
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			err := s.Validate(tt.args.query)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errString) {
				t.Errorf("Set.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSet_NormalizeQuery(t *testing.T) {
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
			want: "SET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet()
			got := s.NormalizeQuery(tt.args.query)
			if got.GetCommand() != tt.want {
				t.Errorf("Set.NormalizeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

package analyzer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/stretchr/testify/mock"
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
			s := NewSet(nil)
			if got := s.Name(); got != tt.want {
				t.Errorf("Set.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Description(t *testing.T) {
	type fields struct {
		Analyzer Analyzer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "check description",
			want: "Set value by key. Usage: SET key value.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet(nil)
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
			s := NewSet(nil)
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
			s := NewSet(nil)
			if got := s.Supports(tt.args.name); got != tt.want {
				t.Errorf("Set.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Validate(t *testing.T) {
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
			errString: "analyzer SET error: cannot process 'ANY' command",
		},
		{
			name:      "fail because wrong argument count",
			args:      args{query: parser.CreateQuery("SET", []string{"Test"})},
			wantErr:   true,
			errString: "analyzer SET error: invalid argument count 2",
		},
		{
			name:      "fail because invalid argument #1",
			args:      args{query: parser.CreateQuery("SET", []string{"==="})},
			wantErr:   true,
			errString: "analyzer SET error: invalid argument #1: ===",
		},
		{
			name:      "fail because invalid argument #2",
			args:      args{query: parser.CreateQuery("SET", []string{"test", "wrong!value"})},
			wantErr:   true,
			errString: "analyzer SET error: invalid argument #2: wrong!value",
		},
		{
			name:      "ok",
			args:      args{query: parser.CreateQuery("SET", []string{"test", "check_value"})},
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet(nil)
			if err := s.Validate(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("Set.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type storageSetMock struct {
	mock.Mock
}

func (m *storageSetMock) Set(key string, value string) error {
	args := m.Called(key, value)

	return args.Error(0)
}

func (m *storageSetMock) Get(key string) (string, bool, error) {
	panic("not implemented") // TODO: Implement
}

func (m *storageSetMock) Del(key string) error {
	panic("not implemented") // TODO: Implement
}

func TestSet_Run(t *testing.T) {
	type args struct {
		query parser.Query
	}

	storageMock1 := new(storageSetMock)
	storageMock1.On("Set", "test", "value").Once().Return(errors.New("test error"))

	storageMock2 := new(storageSetMock)
	storageMock2.On("Set", "test", "value").Once().Return(nil)

	tests := []struct {
		name    string
		storage *storageSetMock
		args    args
		want    any
		wantErr bool
	}{
		{
			name:    "error",
			storage: storageMock1,
			args:    args{query: parser.CreateQuery("TEST", []string{"test", "value"})},
			want:    "",
			wantErr: true,
		},
		{
			name:    "ok",
			storage: storageMock2,
			args:    args{query: parser.CreateQuery("TEST", []string{"test", "value"})},
			want:    "[ok]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet(tt.storage)
			got, err := s.Run(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

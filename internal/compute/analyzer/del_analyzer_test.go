package analyzer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/stretchr/testify/mock"
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
			d := NewDel(nil)
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
			d := NewDel(nil)
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
			d := NewDel(nil)
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
			d := NewDel(nil)
			if got := d.Supports(tt.args.name); got != tt.want {
				t.Errorf("Del.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDel_Validate(t *testing.T) {
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
			errString: "analyzer DEL error: cannot process 'ANY' command",
		},
		{
			name:      "fail because wrong argument count",
			args:      args{query: parser.CreateQuery("DELETE", []string{"Test", "Check"})},
			wantErr:   true,
			errString: "analyzer DEL error: invalid argument count 2",
		},
		{
			name:      "fail because invalid argument",
			args:      args{query: parser.CreateQuery("DEL", []string{"Test&&&?"})},
			wantErr:   true,
			errString: "analyzer DEL error: invalid argument #1: Test&&&?",
		},
		{
			name:      "ok",
			args:      args{query: parser.CreateQuery("DEL", []string{"Test"})},
			wantErr:   false,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel(nil)
			err := d.Validate(tt.args.query)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errString) {
				t.Errorf("Del.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// storage mock
type storageDelMock struct {
	mock.Mock
}

func (m *storageDelMock) Set(_ string, _ string) error {
	panic("not implemented") // TODO: Implement
}

func (m *storageDelMock) Get(_ string) (string, bool, error) {
	panic("not implemented") // TODO: Implement
}

func (m *storageDelMock) Del(key string) error {
	args := m.Called(key)

	return args.Error(0)
}

func TestDel_Run(t *testing.T) {
	type args struct {
		query parser.Query
	}

	storageMock1 := new(storageDelMock)
	storageMock1.On("Del", "test").Once().Return(errors.New("test error"))

	storageMock2 := new(storageDelMock)
	storageMock2.On("Del", "test").Once().Return(nil)

	tests := []struct {
		name    string
		storage *storageDelMock
		args    args
		want    any
		wantErr bool
	}{
		{
			name:    "error",
			storage: storageMock1,
			args:    args{query: parser.CreateQuery("TEST", []string{"test"})},
			want:    "",
			wantErr: true,
		},
		{
			name:    "ok",
			storage: storageMock2,
			args:    args{query: parser.CreateQuery("TEST", []string{"test"})},
			want:    "[ok]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDel(tt.storage)
			got, err := d.Run(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Del.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Del.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

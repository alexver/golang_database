package analyzer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/stretchr/testify/mock"
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
			g := NewGet(nil)
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
			g := NewGet(nil)
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
			g := NewGet(nil)
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
			g := NewGet(nil)
			if got := g.Supports(tt.args.name); got != tt.want {
				t.Errorf("Get.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet_Validate(t *testing.T) {
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
			g := NewGet(nil)
			if err := g.Validate(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("Get.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type storageGetMock struct {
	mock.Mock
}

func (m *storageGetMock) Set(_ string, _ string) error {
	panic("not implemented") // TODO: Implement
}

func (m *storageGetMock) Get(key string) (string, bool, error) {
	args := m.Called(key)

	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *storageGetMock) Del(key string) error {
	panic("not implemented") // TODO: Implement
}

func TestGet_Run(t *testing.T) {
	type args struct {
		query parser.Query
	}

	storageMock1 := new(storageGetMock)
	storageMock1.On("Get", "test").Once().Return("", false, errors.New("test error"))

	storageMock2 := new(storageGetMock)
	storageMock2.On("Get", "test").Once().Return("[not found]", false, nil)

	storageMock3 := new(storageGetMock)
	storageMock3.On("Get", "test").Once().Return("test", true, nil)

	tests := []struct {
		name    string
		storage *storageGetMock
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
			name:    "ok, but not found",
			storage: storageMock2,
			args:    args{query: parser.CreateQuery("TEST", []string{"test"})},
			want:    "[not found]",
			wantErr: false,
		},
		{
			name:    "ok",
			storage: storageMock3,
			args:    args{query: parser.CreateQuery("TEST", []string{"test"})},
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGet(tt.storage)
			got, err := g.Run(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

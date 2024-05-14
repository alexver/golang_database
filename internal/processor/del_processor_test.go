package database

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/stretchr/testify/mock"
)

func TestDelProcessor_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name DEL",
			want: "DEL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDelProcessor(nil)
			if got := p.Name(); got != tt.want {
				t.Errorf("DelProcessor.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelProcessor_Suports(t *testing.T) {
	type args struct {
		query *parser.Query
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no support",
			args: args{query: parser.CreateQuery("TEST", []string{})},
			want: false,
		},
		{
			name: "ok",
			args: args{query: parser.CreateQuery("DEL", []string{})},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDelProcessor(nil)
			if got := p.Suports(tt.args.query); got != tt.want {
				t.Errorf("DelProcessor.Suports() = %v, want %v", got, tt.want)
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

func TestDelProcessor_Process(t *testing.T) {
	type args struct {
		query *parser.Query
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
			p := NewDelProcessor(tt.storage)
			got, err := p.Process(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("DelProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DelProcessor.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

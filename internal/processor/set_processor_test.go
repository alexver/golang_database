package database

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/stretchr/testify/mock"
)

func TestSetProcessor_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name SET",
			want: "SET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSetProcessor(nil)
			if got := p.Name(); got != tt.want {
				t.Errorf("SetProcessor.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetProcessor_Suports(t *testing.T) {
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
			args: args{query: parser.CreateQuery("SET", []string{})},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSetProcessor(nil)
			if got := p.Suports(tt.args.query); got != tt.want {
				t.Errorf("SetProcessor.Suports() = %v, want %v", got, tt.want)
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

func TestSetProcessor_Process(t *testing.T) {
	type args struct {
		query *parser.Query
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
			p := NewSetProcessor(tt.storage)
			got, err := p.Process(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetProcessor.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

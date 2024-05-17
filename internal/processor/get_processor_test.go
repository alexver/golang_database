package database

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/query"
	"github.com/stretchr/testify/mock"
)

func TestGetProcessor_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "check name GET",
			want: "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewGetProcessor(nil)
			if got := p.Name(); got != tt.want {
				t.Errorf("GetProcessor.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProcessor_Suports(t *testing.T) {
	type args struct {
		query *query.Query
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no support",
			args: args{query: query.CreateQuery("TEST", []string{})},
			want: false,
		},
		{
			name: "ok",
			args: args{query: query.CreateQuery("GET", []string{})},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewGetProcessor(nil)
			if got := p.Suports(tt.args.query); got != tt.want {
				t.Errorf("GetProcessor.Suports() = %v, want %v", got, tt.want)
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

func TestGetProcessor_Process(t *testing.T) {
	type args struct {
		query *query.Query
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
			args:    args{query: query.CreateQuery("TEST", []string{"test"})},
			want:    "",
			wantErr: true,
		},
		{
			name:    "ok, but not found",
			storage: storageMock2,
			args:    args{query: query.CreateQuery("TEST", []string{"test"})},
			want:    "[not found]",
			wantErr: false,
		},
		{
			name:    "ok",
			storage: storageMock3,
			args:    args{query: query.CreateQuery("TEST", []string{"test"})},
			want:    "[ok] test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewGetProcessor(tt.storage)
			got, err := p.Process(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProcessor.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

package storage

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alexver/golang_database/internal/storage/engine"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

type InMemoryMock struct {
	mock.Mock
}

func (m *InMemoryMock) Set(key string, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *InMemoryMock) Get(key string) (string, bool, error) {
	args := m.Called(key)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *InMemoryMock) Delete(key string) error {
	args := m.Called(key)

	return args.Error(0)
}

func (m *InMemoryMock) IsSet(string) (bool, error) {
	return false, nil
}

func (m *InMemoryMock) Keys() *[]string {
	var keys = make([]string, 0)

	return &keys
}

func TestCreateStorage(t *testing.T) {
	type args struct {
		engine engine.EngineInterface
		logger *zap.Logger
	}

	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name:    "no engine",
			args:    args{engine: nil},
			wantNil: true,
			wantErr: true,
		},
		{
			name:    "no logger",
			args:    args{engine: new(InMemoryMock), logger: nil},
			wantNil: true,
			wantErr: true,
		},
		{
			name:    "ok",
			args:    args{engine: new(InMemoryMock), logger: zaptest.NewLogger(t)},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateStorage(tt.args.engine, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNil && got != nil {
				t.Errorf("CreateStorage() = %v, want %v", got, tt.wantNil)
			}
			if !tt.wantNil && reflect.TypeOf(got).Kind() != reflect.Pointer {
				t.Errorf("CreateStorage() = %v, want %v", got, reflect.TypeOf(got).Kind())
			}
		})
	}
}

func TestStorage_Set(t *testing.T) {
	type fields struct {
		engine engine.EngineInterface
		logger *zap.Logger
	}
	type args struct {
		key   string
		value string
	}

	var engineMock = new(InMemoryMock)
	engineMock.On("Set", "test_key", "test_value").Once().Return(nil)
	engineMock.On("Set", "test_error", "test_value").Once().Return(errors.New("write lock"))

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{engine: engineMock, logger: zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
				func(e zapcore.Entry) error {
					if e.Level == zap.ErrorLevel {
						t.Error("Should have no error!")
					}
					return nil
				},
			)))},
			args: args{key: "test_key", value: "test_value"},
		},
		{
			name: "error",
			fields: fields{engine: engineMock, logger: zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
				func(e zapcore.Entry) error {
					if e.Level == zap.ErrorLevel {
						if e.Message != "storage SET command error: write lock" {
							t.Errorf("Wrong Error: '%s'", e.Message)
						}
					}
					return nil
				},
			)))},
			args: args{key: "test_error", value: "test_value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				engine: tt.fields.engine,
				logger: tt.fields.logger,
			}
			s.Set(tt.args.key, tt.args.value)
		})
	}
	engineMock.AssertExpectations(t)
}

func TestStorage_Get(t *testing.T) {
	type fields struct {
		engine engine.EngineInterface
		logger *zap.Logger
	}
	type args struct {
		key string
	}

	var engineMock = new(InMemoryMock)
	engineMock.On("Get", "test_key").Once().Return("test_value", true, nil)
	engineMock.On("Get", "test_error").Once().Return("", false, errors.New("read lock"))

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "ok",
			fields: fields{engine: engineMock, logger: zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
				func(e zapcore.Entry) error {
					if e.Level == zap.ErrorLevel {
						t.Error("Should have no error!")
					}
					return nil
				},
			)))},
			args:  args{key: "test_key"},
			want:  "test_value",
			want1: true,
		},
		{
			name: "error",
			fields: fields{engine: engineMock, logger: zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
				func(e zapcore.Entry) error {
					if e.Level == zap.ErrorLevel {
						if e.Message != "storage GET command error: read lock" {
							t.Errorf("Wrong Error: '%s'", e.Message)
						}
					}
					return nil
				},
			)))},
			args:  args{key: "test_error"},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				engine: tt.fields.engine,
				logger: tt.fields.logger,
			}
			got, got1 := s.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Storage.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Storage.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStorage_Del(t *testing.T) {
	type fields struct {
		engine engine.EngineInterface
		logger *zap.Logger
	}
	type args struct {
		key string
	}

	var engineMock = new(InMemoryMock)
	engineMock.On("Delete", "test_key").Once().Return(nil)
	engineMock.On("Delete", "test_error").Once().Return(errors.New("delete lock"))

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{engine: engineMock, logger: zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
				func(e zapcore.Entry) error {
					if e.Level == zap.ErrorLevel {
						t.Error("Should have no error!")
					}
					return nil
				},
			)))},
			args: args{key: "test_key"},
		},
		{
			name: "error",
			fields: fields{engine: engineMock, logger: zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
				func(e zapcore.Entry) error {
					if e.Level == zap.ErrorLevel {
						if e.Message != "storage DEL command error: delete lock" {
							t.Errorf("Wrong Error: '%s'", e.Message)
						}
					}
					return nil
				},
			)))},
			args: args{key: "test_error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				engine: tt.fields.engine,
				logger: tt.fields.logger,
			}
			s.Del(tt.args.key)
		})
	}
}

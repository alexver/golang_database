package engine

import (
	"reflect"
	"sort"
	"testing"
)

var testKeyHash = "92488e1e3eeecdf99f3ed2ce59233efb4b4fb612d5655c0ce9ea52b5a502e655"
var testEmptyKeyHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
var testEmptyValueKeyHash = "275078ecd2550c36aeeb3410fb33eb344c00c8fef187de41991825e1f535114c"

func TestCreateEngine(t *testing.T) {
	tests := []struct {
		name string
		want *dataTable
	}{
		{name: "New Engine", want: &dataTable{data: make(map[string]dataRecord)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateEngine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataTable_Set(t *testing.T) {
	type fields struct {
		data map[string]dataRecord
	}
	type args struct {
		key   string
		value string
	}

	prefilled := make(map[string]dataRecord)
	prefilled[testKeyHash] = dataRecord{key: "test_key", value: "initial value"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		hash    string
		wantErr bool
	}{
		{
			name:    "set value",
			fields:  fields{data: make(map[string]dataRecord)},
			args:    args{key: "test_key", value: "test_value"},
			hash:    testKeyHash,
			wantErr: false,
		},
		{
			name:    "set empty value",
			fields:  fields{data: make(map[string]dataRecord)},
			args:    args{key: "test_key", value: ""},
			hash:    testKeyHash,
			wantErr: false,
		},
		{
			name:    "replace value",
			fields:  fields{data: prefilled},
			args:    args{key: "test_key", value: "new value"},
			hash:    testKeyHash,
			wantErr: false,
		},
		{
			name:    "set value to empty key",
			fields:  fields{data: prefilled},
			args:    args{key: "", value: "test_value"},
			hash:    testEmptyKeyHash,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &dataTable{
				data: tt.fields.data,
			}
			if err := tr.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("dataTable.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tr.data[tt.hash].key != tt.args.key || tr.data[tt.hash].value != tt.args.value {
				t.Errorf("dataTable.Set() record = %v, want record %v", tr.data, tt.args)
			}
		})
	}
}

func Test_dataTable_Get(t *testing.T) {
	type fields struct {
		data map[string]dataRecord
	}
	type args struct {
		key string
	}

	prefilled := preparePrefilledMap()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			name:    "get test_value",
			fields:  fields{data: prefilled},
			args:    args{key: "test_key"},
			want:    "test_value",
			want1:   true,
			wantErr: false,
		},
		{
			name:    "get empty value",
			fields:  fields{data: prefilled},
			args:    args{key: "empty_value_key"},
			want:    "",
			want1:   true,
			wantErr: false,
		},
		{
			name:    "get value by empty key",
			fields:  fields{data: prefilled},
			args:    args{key: ""},
			want:    "empty_key_value",
			want1:   true,
			wantErr: false,
		},
		{
			name:    "key is not found",
			fields:  fields{data: prefilled},
			args:    args{key: "no_key"},
			want:    "",
			want1:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &dataTable{
				data: tt.fields.data,
			}
			got, got1, err := tr.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("dataTable.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dataTable.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("dataTable.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_dataTable_Delete(t *testing.T) {
	type fields struct {
		data map[string]dataRecord
	}
	type args struct {
		key string
	}

	prefilled := preparePrefilledMap()

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		errMessage string
		hash       string
	}{
		{
			name:       "delete test_key",
			fields:     fields{data: prefilled},
			args:       args{key: "test_key"},
			wantErr:    false,
			errMessage: "",
			hash:       testKeyHash,
		},
		{
			name:       "delete empty key",
			fields:     fields{data: prefilled},
			args:       args{key: ""},
			wantErr:    false,
			errMessage: "",
			hash:       testEmptyKeyHash,
		},
		{
			name:       "key is not found",
			fields:     fields{data: prefilled},
			args:       args{key: "no_key"},
			wantErr:    true,
			errMessage: "key no_key is not found",
			hash:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &dataTable{
				data: tt.fields.data,
			}
			err := tr.Delete(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("dataTable.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("dataTable.Delete() error = %v, errMessage %v", err, tt.errMessage)
			}
			if _, ok := tr.data[tt.hash]; ok == true {
				t.Errorf("dataTable.Delete() record = %v, want no key %v", tr.data, tt.args.key)
			}
		})
	}
}

func Test_dataTable_IsSet(t *testing.T) {
	type fields struct {
		data map[string]dataRecord
	}
	type args struct {
		key string
	}

	prefilled := preparePrefilledMap()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "check if exist test_key",
			fields:  fields{data: prefilled},
			args:    args{key: "test_key"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "check if exist empty key",
			fields:  fields{data: prefilled},
			args:    args{key: ""},
			want:    true,
			wantErr: false,
		},
		{
			name:    "key is not found",
			fields:  fields{data: prefilled},
			args:    args{key: "no_key"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &dataTable{
				data: tt.fields.data,
			}
			got, err := tr.IsSet(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("dataTable.IsSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dataTable.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataTable_Keys(t *testing.T) {
	type fields struct {
		data map[string]dataRecord
	}

	prefilled := preparePrefilledMap()
	keys := []string{"", "empty_value_key", "test_key"}

	tests := []struct {
		name   string
		fields fields
		want   *[]string
	}{
		{
			name:   "get all keys",
			fields: fields{data: prefilled},
			want:   &keys,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &dataTable{
				data: tt.fields.data,
			}
			got := tr.Keys()
			sort.Strings(*got)
			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("dataTable.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataTable_getHashAndValue(t *testing.T) {
	type fields struct {
		data map[string]dataRecord
	}
	type args struct {
		key string
	}

	prefilled := preparePrefilledMap()
	prefilled["9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"] = dataRecord{key: "wrong_key_to_test", value: "anything"}
	prefilled["1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014"] = dataRecord{key: "wrong_key_to_test", value: "anything"}
	prefilled["60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752"] = dataRecord{key: "wrong_key_to_test", value: "anything"}
	prefilled["fd61a03af4f77d870fc21e05e7e80678095c92d808cfb3b5c279ee04c74aca13"] = dataRecord{key: "wrong_key_to_test", value: "anything"}
	prefilled["a4e624d686e03ed2767c0abd85c14426b0b1157d2ce81d27bb4fe4f6f01d688a"] = dataRecord{key: "wrong_key_to_test", value: "anything"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		want2   bool
		wantErr bool
	}{
		{
			name:    "key is not found",
			fields:  fields{data: prefilled},
			args:    args{key: "no_key"},
			want:    "5383f79bd56d6f84738075fca2996d3da98050a0775b0cd1a0510817e78f8489",
			want1:   "",
			want2:   false,
			wantErr: false,
		},
		{
			name:    "key is found",
			fields:  fields{data: prefilled},
			args:    args{key: "test_key"},
			want:    testKeyHash,
			want1:   "test_value",
			want2:   true,
			wantErr: false,
		},
		{
			name:    "empty key is found",
			fields:  fields{data: prefilled},
			args:    args{key: ""},
			want:    testEmptyKeyHash,
			want1:   "empty_key_value",
			want2:   true,
			wantErr: false,
		},
		{
			name:    "empty value is found",
			fields:  fields{data: prefilled},
			args:    args{key: "empty_value_key"},
			want:    testEmptyValueKeyHash,
			want1:   "",
			want2:   true,
			wantErr: false,
		},
		{
			name:    "error: too many collisions",
			fields:  fields{data: prefilled},
			args:    args{key: "test"},
			want:    "",
			want1:   "",
			want2:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &dataTable{
				data: tt.fields.data,
			}
			got, got1, got2, err := tr.getHashAndValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("dataTable.getHashAndValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dataTable.getHashAndValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("dataTable.getHashAndValue() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("dataTable.getHashAndValue() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_dataTable_createHash(t *testing.T) {
	type fields struct {
		data map[string]dataRecord
	}
	type args struct {
		key string
		idx int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "key to hash",
			fields: fields{data: make(map[string]dataRecord)},
			args:   args{key: "test", idx: 0},
			want:   "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		{
			name:   "empty key to hash",
			fields: fields{data: make(map[string]dataRecord)},
			args:   args{key: "", idx: 0},
			want:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:   "key to hash with collision",
			fields: fields{data: make(map[string]dataRecord)},
			args:   args{key: "test", idx: 2},
			want:   "60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &dataTable{
				data: tt.fields.data,
			}
			if got := tr.createHash(tt.args.key, tt.args.idx); got != tt.want {
				t.Errorf("dataTable.createHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func preparePrefilledMap() map[string]dataRecord {
	prefilled := make(map[string]dataRecord)
	prefilled[testKeyHash] = dataRecord{key: "test_key", value: "test_value"}
	prefilled[testEmptyKeyHash] = dataRecord{key: "", value: "empty_key_value"}
	prefilled[testEmptyValueKeyHash] = dataRecord{key: "empty_value_key", value: ""}

	return prefilled
}

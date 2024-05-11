package parser

import (
	"reflect"
	"testing"
)

func TestQuery_GetCommand(t *testing.T) {
	tests := []struct {
		name string
		q    *Query
		want string
	}{
		{
			name: "empty query",
			q:    &Query{},
			want: "",
		},
		{
			name: "prefilled",
			q:    &Query{command: "TEST"},
			want: "TEST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.GetCommand(); got != tt.want {
				t.Errorf("Query.GetCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_GetArguments(t *testing.T) {
	tests := []struct {
		name string
		q    *Query
		want []string
	}{
		{
			name: "empty query",
			q:    &Query{},
			want: []string{},
		},
		{
			name: "prefilled",
			q:    &Query{arguments: []string{"", "/*_Test/1290"}},
			want: []string{"", "/*_Test/1290"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.GetArguments(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.GetArguments() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestQuery_GetArgumentCount(t *testing.T) {
	tests := []struct {
		name string
		q    *Query
		want int
	}{
		{
			name: "empty query",
			q:    &Query{},
			want: 0,
		},
		{
			name: "one element",
			q:    &Query{arguments: []string{""}},
			want: 1,
		},
		{
			name: "prefilled few elements",
			q:    &Query{arguments: []string{"", "/*_Test/1290"}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.GetArgumentCount(); got != tt.want {
				t.Errorf("Query.GetArgumentCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

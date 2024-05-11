package parser

import (
	"reflect"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestParser_ParseStringToQuery(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name string
		p    ParserInterface
		args args
		want Query
	}{
		{
			name: "empty string",
			p:    CreatePaser(zaptest.NewLogger(t)),
			args: args{command: ""},
			want: Query{arguments: []string{}},
		},
		{
			name: "command only, no arguments",
			p:    CreatePaser(zaptest.NewLogger(t)),
			args: args{command: "NEW"},
			want: Query{command: "NEW", arguments: []string{}},
		},
		{
			name: "command + 1 argument",
			p:    CreatePaser(zaptest.NewLogger(t)),
			args: args{command: "TEST argument1"},
			want: Query{command: "TEST", arguments: []string{"argument1"}},
		},
		{
			name: "command + 2 arguments",
			p:    CreatePaser(zaptest.NewLogger(t)),
			args: args{command: "TEST argument1 any_CHECK*/"},
			want: Query{command: "TEST", arguments: []string{"argument1", "any_CHECK*/"}},
		},
		{
			name: "command + 3 arguments with spaces",
			p:    CreatePaser(zaptest.NewLogger(t)),
			args: args{command: "TEST\targument1     any_CHECK*/ \t____  \t  \n"},
			want: Query{command: "TEST", arguments: []string{"argument1", "any_CHECK*/", "____"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ParseStringToQuery(tt.args.command); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.ParseStringToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

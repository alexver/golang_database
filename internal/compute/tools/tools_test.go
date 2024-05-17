package tools

import "testing"

func Test_ValidateArgument(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty value",
			args: args{value: ""},
			want: false,
		},
		{
			name: "letter",
			args: args{value: "K"},
			want: true,
		},
		{
			name: "digits",
			args: args{value: "9675"},
			want: true,
		},
		{
			name: "whole set",
			args: args{value: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321/*_"},
			want: true,
		},
		{
			name: "punctuation",
			args: args{value: "!\"#$%%&'()*+,-./:;<=>?@[\\]^_`{|}~"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateArgument(tt.args.value); got != tt.want {
				t.Errorf("ValidateArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

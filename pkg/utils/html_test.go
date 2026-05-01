package utils

import "testing"

func TestRemoveHtmlTags(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test RemoveHtmlTags",
			args: args{
				input: `<div>Hello, world!</div>`,
			},
			want: "Hello, world!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveHtmlTags(tt.args.input); got != tt.want {
				t.Errorf("RemoveHtmlTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

package uuid

import "testing"

func TestNewGoogleUuid(t *testing.T) {
    tests := []struct {
        name string
        want string
    }{
        {
            name: "test google uuid",
            want: "",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Google.Generate()
            t.Logf("NewGoogleUuid = %v", got)
        })
    }
}

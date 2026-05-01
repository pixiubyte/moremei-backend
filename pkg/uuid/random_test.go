package uuid

import "testing"

func Test__RandomClient(t *testing.T) {
    type args struct {
        n int
    }
    tests := []struct {
        name string
        args *args
        want string
    }{
        {
            name: "test random GenerateN",
            args: &args{n: 32},
            want: "",
        },
        {
            name: "test random Generate",
            want: "",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.args != nil {
                t.Logf("GenerateN(%d) = %v", tt.args.n, Random.GenerateN(tt.args.n))
            } else {
                t.Logf("GenerateN() = %v", Random.Generate())
            }
        })
    }
}

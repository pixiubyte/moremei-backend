package jwt

import (
    "fmt"
    "testing"
)

func TestGenerateToken(t *testing.T) {
    type args struct {
        expire int64
        secret string
        params map[string]interface{}
    }
    tests := []struct {
        name    string
        args    args
        want    string
        wantErr bool
    }{
        {
            name: "test generate token",
            args: args{
                expire: 1200,
                secret: "1234567890",
                params: map[string]interface{}{
                    "id": "1234567890",
                },
            },
            want:    "",
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewClient()
            c.SetExpire(tt.args.expire)
            c.SetSecret(tt.args.secret)
            got, err := c.GenerateToken(tt.args.params)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            fmt.Printf("token = %s\n", got)
        })
    }
}

func TestParseToken(t *testing.T) {
    type args struct {
        token  string
        secret string
    }
    tests := []struct {
        name    string
        args    args
        want    map[string]interface{}
        wantErr bool
    }{
        {
            name: "test parse token",
            args: args{
                token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI4ODY2NDAsImlhdCI6MTY3Mjg4NTQ0MCwiaWQiOiIxMjM0NTY3ODkwIn0.DsuVAXg3oM1Y-eI4Mr-dAB-8332ce3JCiYCG-bVfdHA",
                secret: "1234567890",
            },
            want:    nil,
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewClient()
            c.SetSecret(tt.args.secret)
            got, err := c.ParseToken(tt.args.token)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            for k, v := range got {
                fmt.Printf("k = %s, v = %s\n", k, v)
            }
        })
    }
}

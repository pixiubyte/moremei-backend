package cipher

import (
    "reflect"
    "testing"
)

type Fields struct {
    category string
    key      string
    iv       string
}

func TestCipher_Encrypt(t *testing.T) {
    type args struct {
        plaintext string
    }
    tests := []struct {
        name   string
        fields Fields
        args   args
        want   string
    }{
        {
            name: "test encrypt",
            fields: Fields{
                category: "text",
                key:      "1234567890123456",
                iv:       "1234567890123456",
            },
            args: args{
                plaintext: "test",
            },
            want: "WkaOhqSK03Z1pSuPOdc03w==",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c, err := NewCipher(tt.fields.category, tt.fields.key, tt.fields.iv)
            if err != nil {
                t.Errorf("Encrypt() error %v", err)
            }

            encryptResult := c.Encrypt(tt.args.plaintext)

            if got := encryptResult.ToString(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Encrypt() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestCipher_Decrypt(t *testing.T) {
    type args struct {
        ciphertext string
    }
    tests := []struct {
        name   string
        fields Fields
        args   args
        want   string
    }{
        {
            name: "test decrypt",
            fields: Fields{
                category: "text",
                key:      "1234567890123456",
                iv:       "1234567890123456",
            },
            args: args{ciphertext: "WkaOhqSK03Z1pSuPOdc03w=="},
            want: "test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c, err := NewCipher(tt.fields.category, tt.fields.key, tt.fields.iv)
            if err != nil {
                t.Errorf("Decrypt() error %v", err)
                return
            }

            got, err := c.Decrypt(tt.args.ciphertext)
            if err != nil {
                t.Errorf("Decrypt() error %v", err)
                return
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Decrypt() = %v, want %v", got, tt.want)
            }
        })
    }
}

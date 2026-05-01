package uuid

import (
    "time"

    "github.com/zeromicro/go-zero/core/stringx"
)

var Random _RandomClient

type _RandomClient struct{}

func (c *_RandomClient) Generate() string {
    stringx.Seed(time.Now().UnixNano())
    return stringx.Randn(32)
}

func (c *_RandomClient) GenerateN(n int) string {
    stringx.Seed(time.Now().UnixNano())
    return stringx.Randn(n)
}

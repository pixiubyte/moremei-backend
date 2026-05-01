package uuid

import (
    "strings"

    "github.com/zeromicro/go-zero/core/utils"
)

var Google _GoogleClient

type _GoogleClient struct{}

func (u *_GoogleClient) Generate() string {
    return strings.ReplaceAll(utils.NewUuid(), "-", "")
}
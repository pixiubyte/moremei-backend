package uuid

import (
    "github.com/bwmarrin/snowflake"
)

type (
    SnowflakeClient struct {
        Node *snowflake.Node
    }

    SnowflakeConfig struct {
        Node int64 `json:",default=1"`
    }
)

func NewSnowflake(cfg *SnowflakeConfig) *SnowflakeClient {
    newNode, err := snowflake.NewNode(cfg.Node)
    if err != nil {
        panic(err)
    }

    return &SnowflakeClient{Node: newNode}
}

func (s *SnowflakeClient) Generate() string {
    return s.Node.Generate().String()
}

func (s *SnowflakeClient) GenerateInt64() int64 {
    return s.Node.Generate().Int64()
}

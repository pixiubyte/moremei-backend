package uuid

import (
    "testing"
)

func TestSnowflakeClient_Generate(t *testing.T) {
    type fields struct {
        Node int64
    }
    tests := []struct {
        name   string
        fields fields
        want   string
    }{
        {
            name:   "test snowflake uuid",
            fields: fields{Node: 1},
            want:   "",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sf := NewSnowflake(&SnowflakeConfig{Node: tt.fields.Node})
            got := sf.Generate()
            t.Logf("Generate() = %v", got)
        })
    }
}

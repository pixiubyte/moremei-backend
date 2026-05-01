package xftech

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestSkinPutFile(t *testing.T) {
	type fields struct {
		cfg   *Config
		redis redis.RedisConf
	}
	type args struct {
		req *PutFileReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *PutFileResult
		wantErr bool
	}{
		{
			name: "test testskinputfile",
			fields: fields{
				cfg: &Config{
					AppId:  "e6359cea-6567-4dfe-ab00-d8cb572e7d7c",
					Secret: "",
				},
				redis: redis.RedisConf{
					Host: "127.0.0.1:6379",
					Type: "node",
				},
			},
			args: args{
				req: &PutFileReq{
					FaceUrl: "https://b0.bdstatic.com/aac6e932ac6e23b3408f970a87d1d8ff@h_1280",
					Age:     1991,
					Sex:     2,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MustNewClient(tt.fields.cfg, tt.fields.redis)
			got, err := c.TestSkinPutFile(tt.args.req)
			if err != nil {
				t.Errorf("TestSkinPutFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			marshal, err := json.Marshal(*got)
			if err != nil {
				t.Errorf("TestSkinPutFile() error = %v", err)
				return
			}

			fmt.Println(string(marshal))
		})
	}
}

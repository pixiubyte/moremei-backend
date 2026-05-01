package cache

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Redis struct {
	rds *redis.Redis
}

func NewRedis(conf redis.RedisConf) *Redis {
	return &Redis{rds: redis.MustNewRedis(conf)}
}

func (r *Redis) Get(key string) interface{} {
	val, err := r.rds.Get(key)
	if err != nil {
		return nil
	}
	if val == "" {
		return nil
	}
	return val
}

func (r *Redis) Set(key string, val interface{}, timeout time.Duration) error {
	return r.rds.Setex(key, val.(string), int(timeout/time.Second))
}

func (r *Redis) IsExist(key string) bool {
	exists, err := r.rds.Exists(key)
	if err != nil {
		return false
	}
	return exists
}

func (r *Redis) Delete(key string) error {
	_, err := r.rds.Del(key)
	if err != nil {
		return err
	}
	return nil
}

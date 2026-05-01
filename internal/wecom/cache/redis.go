package cache

import (
	"context"
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

func (r *Redis) GetContext(ctx context.Context, key string) interface{} {
	val, err := r.rds.GetCtx(ctx, key)
	if err != nil {
		return nil
	}
	if val == "" {
		return nil
	}
	return val
}

func (r *Redis) SetContext(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return r.rds.SetexCtx(ctx, key, val.(string), int(timeout/time.Second))
}

func (r *Redis) IsExistContext(ctx context.Context, key string) bool {
	exist, err := r.rds.ExistsCtx(ctx, key)
	if err != nil {
		return false
	}
	return exist
}

func (r *Redis) DeleteContext(ctx context.Context, key string) error {
	_, err := r.rds.DelCtx(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

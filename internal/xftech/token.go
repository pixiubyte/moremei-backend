package xftech

import (
	"fmt"
	"sync"
	"time"

	"moremei/ai-saas/internal/xftech/cache"
	"moremei/ai-saas/internal/xftech/util"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type (
	DefaultAccessToken struct {
		conf  *Config
		cache cache.Cache
		lock  *sync.Mutex
	}

	AccessToken struct {
		Token          string `json:"Token"`
		ExpirationTime int    `json:"ExpirationTime"`
	}
)

func NewDefaultAccessToken(conf *Config, cache cache.Cache) (*DefaultAccessToken, error) {
	if cache == nil {
		return nil, errors.New("cache store is nil")
	}

	return &DefaultAccessToken{
		conf:  conf,
		cache: cache,
		lock:  new(sync.Mutex),
	}, nil
}

func (ak *DefaultAccessToken) GetAccessToken() (string, error) {
	// 先从cache中取
	cacheKey := ak.getCacheKey()
	if ak.cache.IsExist(cacheKey) {
		if val := ak.cache.Get(cacheKey); val != nil {
			return val.(string), nil
		}
	}

	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从服务器上获取到不同token
	ak.lock.Lock()
	defer ak.lock.Unlock()

	// 从服务器获取token
	at, err := ak.getTokenFromServer()
	if err != nil {
		return "", err
	}

	// 设置cache
	err = ak.cache.Set(cacheKey, at.Token, time.Duration(at.ExpirationTime-900)*time.Second)
	if err != nil {
		return "", err
	}

	return at.Token, nil
}

func (ak *DefaultAccessToken) getTokenFromServer() (at *AccessToken, err error) {
	commonResponse, err := util.HttpGet(fmt.Sprintf("%s/oauth/token?appid=%s&secret=%s", BaseURL, ak.conf.AppId, ak.conf.Secret))
	if err != nil {
		return nil, err
	}

	if commonResponse.Code != 0 {
		return nil, errors.Errorf("get access_token error : errcode=%v , errormsg=%v", commonResponse.Code, commonResponse.Message)
	}

	var accessToken AccessToken
	err = mapstructure.Decode(commonResponse.Data, &accessToken)
	if err != nil {
		return nil, errors.New("get access_token error : cannot parse data")
	}

	return &accessToken, nil
}

func (ak *DefaultAccessToken) getCacheKey() string {
	return fmt.Sprintf("xftech:%s:token", ak.conf.AppId)
}

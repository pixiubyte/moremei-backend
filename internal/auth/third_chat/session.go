package third_chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/markbates/goth"
)

// Session stores data during the auth process with Wechat.
type Session struct {
	AuthURL      string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	Mobile       string // 用户手机号
}

var _ goth.Session = &Session{}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Wepay provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize the session with Wepay and return the access token to be stored for future use.
func (p *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	// 这里直接返回code作为token
	code := ""
	if params != nil {
		code = params.Get("code")
	} else {
		return "", fmt.Errorf("params is nil")
	}
	p.Mobile = code                                   // code就是手机号
	p.AccessToken = code                              // 直接使用code作为accessToken
	p.ExpiresAt = time.Now().Add(time.Hour * 24 * 14) // 14天过期
	return code, nil
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession wil unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	s := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(s)
	return s, err
}

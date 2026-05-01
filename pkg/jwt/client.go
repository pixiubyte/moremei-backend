package jwt

import (
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/pkg/errors"
)

const (
    jwtAudience  = "aud"
    jwtExpire    = "exp"
    jwtId        = "jti"
    jwtIssueAt   = "iat"
    jwtIssuer    = "iss"
    jwtNotBefore = "nbf"
    jwtSubject   = "sub"
)

type Client struct {
    expire int64
    secret string
}

func NewClient() *Client {
    return &Client{}
}

func (c *Client) SetExpire(expire int64) *Client {
    c.expire = expire
    return c
}

func (c *Client) SetSecret(secret string) *Client {
    c.secret = secret
    return c
}

func (c *Client) GenerateToken(params map[string]interface{}) (string, error) {
    iat := time.Now().Unix()
    claims := make(jwt.MapClaims)
    claims["exp"] = iat + c.expire
    claims["iat"] = iat

    for k, v := range params {
        claims[k] = v
    }

    token := jwt.New(jwt.SigningMethodHS256)
    token.Claims = claims

    return token.SignedString([]byte(c.secret))
}

func (c *Client) ParseToken(token string) (map[string]interface{}, error) {
    t, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(c.secret), nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := t.Claims.(*jwt.MapClaims)
    if !ok || !t.Valid {
        return nil, errors.New("无效的token")
    }

    result := make(map[string]interface{}, len(*claims))
    for k, v := range *claims {
        switch k {
        case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
            continue
        default:
            result[k] = v
        }
    }

    return result, nil
}

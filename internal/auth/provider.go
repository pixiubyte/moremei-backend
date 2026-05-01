package auth

import (
	"fmt"
	"moremei/ai-saas/pkg/utils"
	"net/url"
	"reflect"
	"time"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

type Provider interface {
	GetToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
	goth.Provider
	goth.Session
}

type provider struct {
	goth.Provider
	goth.Session
	config string
	reset  ResetProviderCallback // 重置secret的回调函数
}

func (idp *provider) GetToken(code string) (*oauth2.Token, error) {

	// Need to construct variables supported by goth
	// to call the function to obtain accessToken
	value := url.Values{}
	value.Add("code", code)

	accessToken, err := idp.Session.Authorize(idp.Provider, value)
	if err != nil {
		// 处理Apple授权失败的情况，Apple的secret是动态的，会过期，所以可能需要重新获取
		if idp.Name() == "apple" {
			// Apple auth. It need to request secret dynamically by privateKey.
			if idp.reset != nil {
				err = idp.reset(idp)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
			// try again
			err = nil
			accessToken, err = idp.Session.Authorize(idp.Provider, value)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var expireAt time.Time
	// Get ExpiresAt's value
	valueOfExpire := reflect.ValueOf(idp.Session).Elem().FieldByName("ExpiresAt")
	if valueOfExpire.IsValid() {
		expireAt = valueOfExpire.Interface().(time.Time)
	}
	token := oauth2.Token{
		AccessToken: accessToken,
		Expiry:      expireAt,
	}

	return &token, nil
}

// 需要先调用 GetToken 方法获取token
func (idp *provider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	gothUser, err := idp.Provider.FetchUser(idp.Session)
	if err != nil {
		return nil, err
	}
	return getUser(gothUser, idp.Provider.Name()), nil
}

func getUser(gothUser goth.User, provider string) *UserInfo {
	user := UserInfo{
		Id:       gothUser.UserID,
		Name:     gothUser.Name,
		NickName: gothUser.NickName,
		Email:    gothUser.Email,
		Avatar:   gothUser.AvatarURL,
	}
	// Some idp return an empty Name
	// so construct the Name with firstname and lastname or nickname
	if user.Name == "" {
		if gothUser.FirstName != "" && gothUser.LastName != "" {
			user.Name = getName(gothUser.FirstName, gothUser.LastName)
		} else {
			user.Name = gothUser.NickName
		}
	}
	if user.NickName == "" {
		if gothUser.FirstName != "" && gothUser.LastName != "" {
			user.NickName = getName(gothUser.FirstName, gothUser.LastName)
		} else {
			user.NickName = user.Name
		}
	}
	if provider == "steam" {
		user.Name = user.Id
		user.Email = ""
	} else if provider == "apple" {
		user.Name = utils.GetUsernameFromEmail(user.Email)
	}
	if gothUser.RawData != nil {
		uniid := gothUser.RawData["Unionid"]
		if uid, ok := uniid.(string); ok {
			user.UnionId = uid
		}
		user.Extra = gothUser.RawData
	}
	if user.UnionId == "" { // 无unionid时，使用id作为unionid
		user.UnionId = user.Id
	}
	return &user
}

func getName(firstName, lastName string) string {
	if utils.IsChinese(firstName) || utils.IsChinese(lastName) {
		return fmt.Sprintf("%s%s", lastName, firstName)
	} else {
		return fmt.Sprintf("%s %s", firstName, lastName)
	}
}

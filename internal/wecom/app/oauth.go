package app

import (
	"fmt"

	"moremei/ai-saas/internal/wecom/util"
)

const (
	getOauthUserInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo"
)

type GetOauthUserInfoRequest struct {
	Code string `json:"code"`
}
type GetOauthUserInfoResponse struct {
	util.CommonError
	UserId         string `json:"userid,omitempty"`
	ExternalUserId string `json:"external_userid,omitempty"`
	ParentUserId   string `json:"parent_userid,omitempty"`
	Openid         string `json:"openid,omitempty"`
}

func (r *Client) GetOauthUserInfo(request GetOauthUserInfoRequest) (*GetOauthUserInfoResponse, error) {
	accessToken, err := r.GetAccessToken()
	if err != nil {
		return nil, err
	}
	var response []byte
	response, err = util.HTTPGet(fmt.Sprintf("%s?access_token=%v&code=%v", getOauthUserInfoURL, accessToken, request.Code))
	if err != nil {
		return nil, err
	}
	var result GetOauthUserInfoResponse
	err = util.DecodeWithError(response, &result, "GetOauthUserInfoResponse")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

package xftech

import (
	"fmt"
	"strconv"

	"moremei/ai-saas/internal/xftech/errox"
	"moremei/ai-saas/internal/xftech/util"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type (
	PutFile struct {
		baseUrl string
		ak      *DefaultAccessToken
	}

	PutFileReq struct {
		FaceUrl string
		Age     int64
		Sex     int64
	}

	PutFileResult struct {
		DistinguishId   string `json:"DistinguishId"`
		ImgUrl          string `json:"ImgUrl"`
		RedUrl          string `json:"RedUrl"`
		BrownUrl        string `json:"BrownUrl"`
		EnhanceUrl      string `json:"EnhanceUrl"`
		TImgUrl         string `json:"TImgUrl"`
		LeftImgUrl      string `json:"LeftImgUrl"`
		RightImgUrl     string `json:"RightImgUrl"`
		ChinImgUrl      string `json:"ChinImgUrl"`
		PointUrl        string `json:"PointUrl"`
		Sex             int64  `json:"Sex"`
		ColorLevel      uint64 `json:"ColorLevel"`
		HueLevel        uint64 `json:"HueLevel"`
		Skin            uint64 `json:"Skin"`
		T               uint64 `json:"t"`
		Cheek           uint64 `json:"Cheek"`
		Chin            uint64 `json:"Chin"`
		Pore            uint64 `json:"Pore"`
		PoreNum         uint64 `json:"PoreNum"`
		Blackhead       uint64 `json:"Blackhead"`
		BlackheadNum    uint64 `json:"BlackheadNum"`
		Acne            uint64 `json:"Acne"`
		AcneNum         uint64 `json:"AcneNum"`
		DarkSpots       uint64 `json:"DarkSpots"`
		DarkSpotsNum    uint64 `json:"DarkSpotsNum"`
		AcneScarringNum uint64 `json:"AcneScarringNum"`
		NevusNum        uint64 `json:"NevusNum"`
		BlackEye        uint64 `json:"BlackEye"`
		LeftBlackEye    uint64 `json:"LeftBlackEye"`
		RightBlackEye   uint64 `json:"RightBlackEye"`
		Eyebag          uint64 `json:"Eyebag"`
		LeftEyebag      uint64 `json:"LeftEyebag"`
		RightEyebag     uint64 `json:"RightEyebag"`
		Wrinkle         uint64 `json:"Wrinkle"`
		RaiseHead       uint64 `json:"RaiseHead"`
		DecreePattern   uint64 `json:"DecreePattern"`
		LineWrinkles    uint64 `json:"LineWrinkles"`
		Sensitive       uint64 `json:"Sensitive"`
		Moisture        uint64 `json:"Moisture"`
		Vascular        uint64 `json:"Vascular"`
		Pigment         uint64 `json:"Pigment"`
		Structure       uint64 `json:"Structure"`
		RedLevel        uint64 `json:"RedLevel"`
		BrownLevel      uint64 `json:"BrownLevel"`
	}
)

func (c *Client) TestSkinPutFile(req *PutFileReq) (*PutFileResult, error) {
	accessToken, err := c.ak.GetAccessToken()
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"token":   accessToken,
		"faceUrl": req.FaceUrl,
		"sex":     strconv.FormatInt(req.Sex, 10),
		"age":     strconv.FormatInt(req.Age, 10),
	}

	url := fmt.Sprintf("%s/testskin/putfile", BaseURL)
	commonResponse, err := util.HttpPostWithFormData(url, data)
	if err != nil {
		return nil, err
	}

	if commonResponse.Code != CodeSuccess {
		if commonResponse.Code >= 103 && commonResponse.Code <= 109 {
			return nil, errox.NewBizError(commonResponse.Code, commonResponse.Message)
		}
		return nil, errors.Errorf("%s error : errcode=%v , errormsg=%v", url, commonResponse.Code, commonResponse.Message)
	}

	var putFileResult PutFileResult
	err = mapstructure.Decode(commonResponse.Data, &putFileResult)
	if err != nil {
		return nil, errors.Errorf("%s error  : cannot parse data", url)
	}

	return &putFileResult, nil
}

package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/uuid"
	errorx "moremei/ai-saas/pkg/x/error"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var filePutAllowedExt = map[string][]string{
	consts.FilePutSignCategorySkinTest:      {"png", "jpg", "jpeg"},
	consts.FilePutSignCategoryAvatar:        {"png", "jpg", "jpeg"},
	consts.FilePutSignCategoryUserPhoto:     {"png", "jpg", "jpeg"},
	consts.FilePutSignCategoryMedicalReport: {"png", "jpg", "jpeg", "pdf", "doc", "docx"},
}

type FilePutSignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type FilePutKey struct {
	Category string
	UserUuid string
	Filename string
}

func NewFilePutSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilePutSignLogic {
	return &FilePutSignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FilePutSignLogic) FilePutSign(req *types.FilePutSignRequest) (resp *types.FilePutSignResponse, err error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		return nil, errorx.UnauthorizedError
	}

	fileExt := strings.ToLower(req.FileExt)
	if !stringx.Contains(filePutAllowedExt[req.Category], fileExt) {
		return nil, errs.CosFileCategoryError
	}

	filePutKey := NewFilePutKey(req.Category, user.Uuid, fileExt)
	putKey := filePutKey.ToString()

	var opt sts.CredentialOptions

	switch req.Category {
	case consts.FilePutSignCategorySkinTest, consts.FilePutSignCategoryAvatar,
		consts.FilePutSignCategoryUserPhoto, consts.FilePutSignCategoryMedicalReport:
		opt = sts.CredentialOptions{
			Policy: &sts.CredentialPolicy{
				Statement: []sts.CredentialPolicyStatement{
					{
						Action: []string{
							"name/cos:PostObject",
							"name/cos:PutObject",
						},
						Effect: "allow",
						Resource: []string{
							fmt.Sprintf(
								"qcs::cos:%s:uid/%s:%s/%s",
								l.svcCtx.Config.Cos.Region,
								l.svcCtx.Config.Cos.AppId,
								l.svcCtx.Config.Cos.Bucket,
								putKey,
							),
						},
					},
				},
			},
			Region:          l.svcCtx.Config.Cos.Region,
			DurationSeconds: int64(time.Hour.Seconds()),
		}
	default:
		return nil, errorx.InvalidParamsError
	}

	cosStsClient := l.svcCtx.CosClient.GetStsClient()
	cosClient := l.svcCtx.CosClient.GetCosClient()

	credential, err := cosStsClient.GetCredential(&opt)
	if err != nil {
		return nil, err
	}

	authorization := cosClient.Object.GetSignature(
		l.ctx,
		http.MethodPut,
		putKey,
		credential.Credentials.TmpSecretID,
		credential.Credentials.TmpSecretKey,
		time.Hour,
		nil,
	)

	return &types.FilePutSignResponse{
		SecurityToken: credential.Credentials.SessionToken,
		Authorization: authorization,
		PutUrl:        cosClient.Object.GetObjectURL(putKey).String(),
		PutHost:       cosClient.BaseURL.BucketURL.String(),
		PutKey:        putKey,
	}, nil
}

func NewFilePutKey(category, userUuid, fileExt string) *FilePutKey {
	return &FilePutKey{
		Category: category,
		UserUuid: userUuid,
		Filename: fmt.Sprintf("%s%s.%s", uuid.Google.Generate(), time.Now().Format("20060102"), fileExt),
	}
}

func ParseToFilePutKey(category, userUuid, putKey string) (*FilePutKey, error) {
	data := strings.Split(putKey, "/")
	if len(data) != 3 {
		return nil, errorx.InvalidParamsError
	}
	if data[0] != category {
		return nil, errs.CosFileForbiddenError
	}
	_, ok := filePutAllowedExt[data[0]]
	if !ok {
		return nil, errs.CosFileForbiddenError
	}
	if data[1] != userUuid {
		return nil, errs.CosFileForbiddenError
	}

	return &FilePutKey{
		Category: data[0],
		UserUuid: data[1],
		Filename: data[2],
	}, nil
}

func (fpk *FilePutKey) ToString() string {
	return fmt.Sprintf("%s/%s/%s", fpk.Category, fpk.UserUuid, fpk.Filename)
}

package captcha

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/threading"
	"golang.org/x/net/context"
)

type Client struct {
	secretId         string
	secretKey        string
	region           string
	captchaAppId     uint64
	captchaAppSecret string
	expire           int64
	rds              *redis.Redis
}

type Scene string

const (
	SceneMiniProgram Scene = "mini"
	SceneWeb         Scene = "web"
)

func (s Scene) String() string {
	return string(s)
}

type VerifyOptions struct {
	ticketCacheKeyPrefix string

	ticket  *string
	userIp  *string
	randStr *string
}

type VerifyOption func(*VerifyOptions)

func WithTicket(ticket string) VerifyOption {
	return func(o *VerifyOptions) {
		o.ticket = &ticket
	}
}
func WithUserIp(userIp string) VerifyOption {
	return func(o *VerifyOptions) {
		o.userIp = &userIp
	}
}
func WithRandStr(randStr string) VerifyOption {
	return func(o *VerifyOptions) {
		o.randStr = &randStr
	}
}
func WithTicketCacheKeyPrefix(prefix string) VerifyOption {
	return func(o *VerifyOptions) {
		o.ticketCacheKeyPrefix = prefix
	}
}

func NewClient(cfg Config, rds *redis.Redis) *Client {
	return &Client{
		secretId:         cfg.SecretId,
		secretKey:        cfg.SecretKey,
		region:           cfg.Region,
		captchaAppId:     cfg.CaptchaAppId,
		captchaAppSecret: cfg.CaptchaAppSecret,
		expire:           cfg.Expire,
		rds:              rds,
	}
}

// EncryptCaptchaAppId 加密 CaptchaAppid, see https://cloud.tencent.com/document/product/1110/36841?from=console_document_search#c6a3e517-4962-4932-9bf8-9fb10fd03166
// 加密规则：
// 1. 选择所需 CaptchaAppSecretKey，作为密钥 key。当 key 小于32字节时，循环填充同一个 CaptchaAppSecretKey 补齐到32字节作为密钥 key。
// 2. 随机生成16字节的 IV，结合步骤1中得到的密钥 Key，对业务数据对象进行 AES256加密，加密模式为 CBC/PKCS7Padding，
//    加密的业务数据对象为验证码业务 CaptchaAppid&时间戳&密文过期时间，时间戳为秒级当前 unix 时间戳，不可为未来时间，
//    密文过期时间单位为秒最大值为86400秒，得到加密后的串为 CaptchaAppidEncrypted。
// 3. 对步骤2中 IV 拼接 AES256加密得到的字节数组 CaptchaAppidEncrypted 进行 Base64编码，即 Base64（IV+CaptchaAppidEncrypted），
//    中间无连接符，得到最终加密后的业务参数请求字符串即为 aidEncrypted。
func (c *Client) EncryptCaptchaAppId() (string, error) {
	// 补充密钥到 32 位
	appSecretKey := []byte(c.captchaAppSecret)
	remainder := 32 % len(appSecretKey)
	key := append(appSecretKey, appSecretKey[:remainder]...)

	// 拼接待加密的明文
	curTime := time.Now().Unix()
	plaintext := strconv.FormatUint(c.captchaAppId, 10) + "&" + strconv.FormatInt(curTime, 10) + "&" +
		strconv.FormatInt(c.expire, 10)

	// 生成 16 字节的 IV
	iv := []byte(stringx.Randn(16))

	// 创建 AES 块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充明文
	plaintextPadded := c.pkcs7Pad([]byte(plaintext), block.BlockSize())

	// 创建 CBC 模式加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密
	ciphertext := make([]byte, len(plaintextPadded))
	mode.CryptBlocks(ciphertext, plaintextPadded)

	// 拼接 IV 和密文并进行 Base64 编码
	ivAndCiphertext := append(iv, ciphertext...)

	return base64.StdEncoding.EncodeToString(ivAndCiphertext), nil
}

func (c *Client) VerifyTicket(scene Scene, options ...VerifyOption) (bool, error) {
	opt := &VerifyOptions{}
	for _, o := range options {
		o(opt)
	}

	if scene == SceneMiniProgram {
		if opt.ticket == nil || opt.userIp == nil {
			return false, fmt.Errorf("小程序验证码校验缺少必要参数")
		}
	} else if scene == SceneWeb {
		if opt.ticket == nil || opt.userIp == nil || opt.randStr == nil {
			return false, fmt.Errorf("网页验证码校验缺少必要参数")
		}
	} else {
		return false, fmt.Errorf("不支持的验证码校验来源：%scene", scene)
	}

	ticketCacheKey := opt.ticketCacheKeyPrefix + "captcha:ticket:" + hash.Md5Hex([]byte(*opt.ticket))
	exists, err := c.rds.ExistsCtx(context.Background(), ticketCacheKey)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	defer func() {
		threading.GoSafe(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			if err := c.rds.SetexCtx(ctx, ticketCacheKey, *opt.ticket, 60*60*24); err != nil {
				logx.Errorf("设置验证码票据 `%s` 缓存失败：%v", *opt.ticket, err)
			}
		})
	}()

	credential := common.NewCredential(c.secretId, c.secretKey)

	prof := profile.NewClientProfile()
	prof.NetworkFailureMaxRetries = 3                             // 定义最大重试次数
	prof.NetworkFailureRetryDuration = profile.ExponentialBackoff // 定义重试间隔时间

	client, _ := captcha.NewClient(credential, c.region, prof)

	var (
		code int64
		msg  string
	)

	if scene == SceneMiniProgram {
		request := captcha.NewDescribeCaptchaMiniResultRequest()
		request.CaptchaType = common.Uint64Ptr(9)
		request.Ticket = opt.ticket
		request.UserIp = opt.userIp
		request.CaptchaAppId = &c.captchaAppId
		request.AppSecretKey = &c.captchaAppSecret
		response, err := client.DescribeCaptchaMiniResult(request)
		if err != nil {
			return false, fmt.Errorf("腾讯云验证码校验失败：%w", err)
		}

		if response.Response == nil || response.Response.CaptchaCode == nil {
			return false, fmt.Errorf("腾讯云验证码校验失败：%s", response.ToJsonString())
		}

		code = *response.Response.CaptchaCode
		msg = *response.Response.CaptchaMsg
	} else {
		request := captcha.NewDescribeCaptchaResultRequest()
		request.CaptchaType = common.Uint64Ptr(9)
		request.Ticket = opt.ticket
		request.UserIp = opt.userIp
		request.Randstr = opt.randStr
		request.CaptchaAppId = &c.captchaAppId
		request.AppSecretKey = &c.captchaAppSecret

		response, err := client.DescribeCaptchaResult(request)
		if err != nil {
			return false, fmt.Errorf("腾讯云验证码校验失败：%w", err)
		}

		if response.Response == nil || response.Response.CaptchaCode == nil {
			return false, fmt.Errorf("腾讯云验证码校验失败：%s", response.ToJsonString())
		}

		// 无感模式下，需要根据返回的 EvilLevel 进行处理
		if response.Response.EvilLevel != nil && *response.Response.EvilLevel == 100 {
			// TODO:需要返回一个错误信息来提示用户重新输入验证码
			return false, nil
		}

		code = *response.Response.CaptchaCode
		msg = *response.Response.CaptchaMsg
	}

	// see https://cloud.tencent.com/document/product/1110/48499#3.-.E8.BE.93.E5.87.BA.E5.8F.82.E6.95.B0
	switch code {
	case 1:
		// 验证成功
		return true, nil
	case 7, 10, 15, 16, 21, 25, 100:
		// 无效票据
		return false, nil
	default:
		return false, fmt.Errorf("腾讯云验证码校验失败：%s", msg)
	}
}

// pkcs7Pad 对数据进行 PKCS7 填充
func (c *Client) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

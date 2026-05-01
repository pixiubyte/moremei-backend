package svc

import (
	"moremei/ai-saas/cmd/tenant/api/internal/config"
	"moremei/ai-saas/cmd/tenant/api/internal/middleware"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/cos"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/uuid"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Snowflake *uuid.SnowflakeClient
	DB        *gorm.DB

	RedisClient *redis.Redis
	CosClient   *cos.Client

	Config config.Config

	UserModel              tenant.UserModel
	WebUserModel           tenant.WebUserModel
	AuthProviderConfModel  tenant.AuthProviderConfModel
	AuthProviderUserModel  tenant.AuthProviderUserModel
	AiRobotModel           tenant.AiRobotModel
	AiRobotWhitelistModel  tenant.AiRobotWhitelistModel
	ConversationModel      tenant.ConversationModel
	MessageModel           tenant.MessageModel
	MessageOpinionModel    tenant.MessageOpinionModel
	SharedChatHistoryModel tenant.SharedChatHistoryModel
	SmsModel               tenant.SmsModel
	MiniUserModel          tenant.MiniUserModel
	SystemConfigModel      tenant.SystemConfigModel
	ParamsStorageModel     tenant.ParamsStorageModel
	AppUserProfileModel    tenant.AppUserProfileModel

	TokenAuthMiddleware      rest.Middleware
	TokenManageModel         tenant.TokenManageModel
	TokenRequestHistoryModel tenant.TokenRequestHistoryModel
	TokenCheckMiddleware     rest.Middleware

	SdkClients *SdkClients
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gormc.NewClient(&c.MySQL)
	if err != nil {
		panic(err)
	}
	redisClient := redis.MustNewRedis(c.Redis)
	if !redisClient.Ping() {
		panic("redis's connection is invalid")
	}

	systemConfigModel := tenant.NewSystemConfigModel(db)
	tokenManageModel := tenant.NewTokenManageModel(db)
	tokenRequestHistoryModel := tenant.NewTokenRequestHistoryModel(db)

	return &ServiceContext{
		Snowflake: uuid.NewSnowflake(&c.Snowflake),
		DB:        db,

		RedisClient: redisClient,
		CosClient:   cos.NewClient(c.Cos),

		Config: c,

		UserModel:              tenant.NewUserModel(db),
		WebUserModel:           tenant.NewWebUserModel(db),
		AuthProviderConfModel:  tenant.NewAuthProviderConfModel(db),
		AuthProviderUserModel:  tenant.NewAuthProviderUserModel(db),
		AiRobotModel:           tenant.NewAiRobotModel(db),
		AiRobotWhitelistModel:  tenant.NewAiRobotWhitelistModel(db),
		ConversationModel:      tenant.NewConversationModel(db),
		MessageModel:           tenant.NewMessageModel(db),
		MessageOpinionModel:    tenant.NewMessageOpinionModel(db),
		SharedChatHistoryModel: tenant.NewSharedChatHistoryModel(db),
		SmsModel:               tenant.NewSmsModel(db),
		MiniUserModel:          tenant.NewMiniUserModel(db),
		SystemConfigModel:      systemConfigModel,
		ParamsStorageModel:     tenant.NewParamsStorageModel(db),
		AppUserProfileModel:    tenant.NewAppUserProfileModel(db),

		TokenAuthMiddleware:      middleware.NewTokenAuthMiddleware(c.Name, tokenManageModel, tokenRequestHistoryModel, redisClient).Handle,
		TokenManageModel:         tokenManageModel,
		TokenRequestHistoryModel: tokenRequestHistoryModel,
		TokenCheckMiddleware:     middleware.NewTokenCheckMiddleware(&c).Handle,

		SdkClients: NewSdkClients(c, systemConfigModel, redisClient),
	}
}

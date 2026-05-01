package consts

type CtxKey string

const (
	UserIDCtxKey   = "au_uid"
	RemoteIpCtxKey = "remote_ip"
	MeiAngelCtxKey = "mei_angel"

	SvcAdvisorCtxKey CtxKey = "svc_advisor" // 服务商

	SmsCodeLength              = 6
	SmsCodeExpireSeconds       = 300
	SmsCodeRepeatExpireSeconds = 60

	MessageOpinionActionThumbsUp   = "thumbsup"
	MessageOpinionActionThumbsDown = "thumbsdown"
	MessageOpinionActionCancel     = "cancel"

	SubNameConversation = "conversation_sub"

	FilePutSignCategorySkinTest      = "skintest"
	FilePutSignCategoryAvatar        = "avatar"
	FilePutSignCategoryUserPhoto     = "userphoto"
	FilePutSignCategoryMedicalReport = "medicalreport"

	RecommendQuestionsCacheKey = "recommend:questions"

	QrcodeList = "%s:%s:qrcode:list"

	MiniProgramTokenManage = "%s:token:manage:%s"

	QuestionnaireLimit = "%s:%s:questionnaire:limit:%s"

	// 文章阅读数和点赞数缓存key
	ArticleViewCountKey    = "article:views"
	ArticleLikeCountKey    = "article:likes"
	ArticleCollectCountKey = "article:collects"
	ArticleSharedCountKey  = "article:shared"

	// 小程序分享小程序码存储分类
	FilePutSignCategoryMiniShareQRCode = "minishare"

	// 任务处理key
	UserMissionKey = "user:mission"

	// 天气缓存
	WeatherCacheKey = "weather"

	AppNameDefault = ""

	UserProfileKeyArchive = "user:profile:health:archive"

	BloodSugarType = "blood-sugar"

	CheckInTypeClick = 1
	CheckInTypePhoto = 2
	StatusCheckedIn  = 1 // 已打卡
	StatusNotChecked = 2 // 未打卡
	CycleTypeDay     = 1
	CycleTypeCustom  = 2

	MissionTypeDaily  = 0
	MissionTypeSingle = 1
	MissionTypeAll    = 2

	MissionsPartiallyChecked = 2
	MissionsStatusNotChecked = 0
	MissionsStatusChecked    = 1

	BloodPressureMetrics = "blood_pressure"
	BloodGlucoseMetrics  = "blood_glucose"
	UricAcidMetrics      = "uric_acid"
	BloodLipidsMetrics   = "blood_lipids"
	WeightMetrics        = "weight"

	UserArchiveStatusActive = 2
	UserProfileCacheKey     = "%s:user:profile:%d"

	HabitDailyTipKey = "%s:daily_tip:answer_count:%d:%s"

	UserMissionHomeQuestionDetail  = "问一问"
	UserMissionDailyShareDetail    = "分享文章"
	UserMissionPhoneRegisterDetail = "手机注册"
	UserMissionBasicInfoDetail     = "完成基本信息填写"

	UserPointsRecordType = "mission"

	// 疾病类型
	DiseaseDiabetes       = "糖尿病"
	DiseaseHypertension   = "高血压"
	DiseaseHyperuricemia  = "痛风"
	DiseaseHyperlipidemia = "高血脂"

	// 指标类型
	MetricTypeBloodPressure = "blood_pressure"
	MetricTypeBloodGlucose  = "blood_glucose"
	MetricTypeUricAcid      = "uric_acid"
	MetricTypeBloodLipids   = "blood_lipids"
	MetricTypeWeight        = "weight"

	// 测量状态
	MeasureStatusFasting   = "空腹"
	MeasureStatusAfterMeal = "餐后"
	MeasureStatusRandom    = "随机"

	DailyTipsKey          = "%s:daily:tip:%s"
	QuestionnaireItemsKey = "%s:questionnaire_item_ids:%d:%d"

	// 健康组通知类型常量
	HealthGroupNotificationTypeJoinSuccess  = "group_join_success"  // 加入成功通知
	HealthGroupNotificationTypeMemberJoined = "group_member_joined" // 新成员加入通知
	HealthGroupNotificationTypeJoinRequest  = "group_join_request"  // 加入申请通知(管理员接收)
	HealthGroupNotificationTypeJoinApplied  = "group_join_applied"  // 申请已提交通知
	HealthGroupNotificationTypeJoinResult   = "group_join_result"   // 申请结果通知
	HealthGroupNotificationTypeQuitResult   = "group_quit_result"   // 退出群组通知

	// 健康组通知标题常量
	HealthGroupNotificationTitleJoinSuccess  = "加入家庭组成功"
	HealthGroupNotificationTitleMemberJoined = "新成员加入家庭组"
	HealthGroupNotificationTitleJoinRequest  = "新的加入申请"
	HealthGroupNotificationTitleJoinApplied  = "加入申请已提交"
	HealthGroupNotificationTitleJoinResult   = "家庭组申请结果"
	HealthGroupNotificationTitleQuit         = "用户退出家庭组"

	// 健康组通知目标类型
	HealthGroupNotificationTargetTypeGroup = "health_group"

	// 健康组通知状态
	HealthGroupNotificationStatusUnread = 0 // 未读

	// 健康指标通知类型常量
	HealthGroupNotificationTypeMetricsAlert = "health_metrics_alert" // 健康指标提醒类型

	// 健康指标通知目标类型
	HealthGroupNotificationTargetTypeMetrics = "health_metrics" // 健康指标目标类型
	HealthGroupNotificationTargetTypeHabit   = "habit"          // 健康习惯目标类型

	// 健康指标异常类型
	HealthMetricsAbnormalTypeHigh = "high" // 指标偏高
	HealthMetricsAbnormalTypeLow  = "low"  // 指标偏低

	// 健康指标通知标题常量
	HealthGroupNotificationTitle = "健康指标异常提醒"

	// 习惯打卡通知相关常量
	HealthGroupNotificationTypeHabitReminder  = "habit_reminder" // 习惯打卡提醒类型
	HealthGroupNotificationTitleHabitRemind   = "习惯打卡提醒"         // 习惯打卡提醒标题
	HealthGroupNotificationContentHabitRemind = "您的习惯「%s」该打卡啦"   // 健康组通知内容常量

	// 健康组通知内容常量
	HealthGroupNotificationContentJoinSuccess  = "加入家庭组「%s」"
	HealthGroupNotificationContentMemberJoined = "用户「%s」已加入家庭组「%s」"
	HealthGroupNotificationContentMemberQuit   = "用户「%s」已退出家庭组「%s」"
	HealthGroupNotificationContentJoinRequest  = "用户「%s」申请加入家庭组「%s」"
	HealthGroupNotificationContentJoinApplied  = "申请加入家庭组「%s」"

	// 健康组通知简短内容常量
	HealthGroupNotificationShortContentApproved    = "审核通过"
	HealthGroupNotificationShortContentNewMember   = "新成员加入"
	HealthGroupNotificationShortContentJoinRequest = "申请加入"
	HealthGroupNotificationShortContentApplied     = "申请已提交"
	HealthGroupNotificationShortContentReject      = "申请已拒绝" // 新增：申请拒绝的简短内容
	HealthGroupNotificationShortContentApprove     = "申请已通过" // 新增：申请通过的简短内容
	HealthGroupNotificationShortContentQuit        = "退出家庭组" // 新增：申请通过的简短内容

	// 健康组通知内容常量
	HealthGroupNotificationContentJoinResult        = "加入家庭组「%s」的申请%s" // 新增：申请结果通知内容
	HealthGroupNotificationContentJoinResultReject  = "已被拒绝"           // 新增：申请被拒绝的结果文本
	HealthGroupNotificationContentJoinResultApprove = "已通过"            // 新增：申请通过的结果文本

	HabitNotificationExtraDataHabitId      = "habit_id"      // 习惯ID字段
	HabitNotificationExtraDataHabitName    = "habit_name"    // 习惯名称字段
	HabitNotificationExtraDataReminderTime = "reminder_time" // 提醒时间字段

	QuestionnaireSceneChat    = "chat"
	QuestionnaireSceneDefault = "default"

	ModuleUserReportHeader           = "user_report_header"
	ModuleHealthAssessment           = "health_assessment"
	ModulePeerComparison             = "peer_comparison"
	ModuleMainIssuesSymptoms         = "main_issues_symptoms"
	ModuleHealthGoals                = "health_goals"
	ModuleProfessionalRecommendation = "professional_recommendation"
	ModuleLifestyleHealthPlan        = "lifestyle_health_plan"
	ModuleProductInfo                = "productinfo"
	ModuleHealthPlanTasks            = "health_plan_tasks"
	ModuleProductEfficacy            = "product_efficacy"
	ModuleHealthIndicators           = "health_indicators"
	ModulePersonalizedHealthPlan     = "personalized_health_plan"

	FormTypeFixed   = 1 // 固定问卷
	FormTypeDynamic = 2 // 动态问卷

	ArticleRecommendsListKey = "%s:recommends:article:%s"
	SceneMultiModal          = "multimodal"
)

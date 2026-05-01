package tenant

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"

	"gorm.io/gorm"
)

const (
	SystemConfigGroupApp         = "app"
	SystemConfigGroupAppointment = "appointment"
	SystemConfigGroupMember      = "member"
	SystemConfigGroupHezhen      = "hezhen"
	SystemConfigGroupRmy         = "rmy"
	SystemConfigGroupMeiManager  = "mei_manager"

	SystemConfigKeyAppName                  = "app_name"
	SystemConfigKeyAppLogo                  = "app_logo"
	SystemConfigKeyAppTheme                 = "app_theme"
	SystemConfigKeyAppLanguage              = "app_language"
	SystemConfigKeyAppSideCollapsed         = "app_side_collapsed"
	SystemConfigKeyAppUserAgreement         = "app_user_agreement" // 用户协议
	SystemConfigKeyAppPrivacyPolicy         = "app_privacy_policy" // 隐私政策
	SystemConfigKeyAppTechnicalSupport      = "app_technical_support"
	SystemConfigKeyAppWeatherIcons          = "app_weather_icons"
	SystemConfigKeyMessageThemeAbout        = "theme_about"
	SystemConfigKeyMessageFirstMessage      = "first_message"
	SystemConfigKeyMessageSkinTestUrl       = "skin_test_url"
	SystemConfigKeyBloodPressureNote        = "blood_pressure_note"
	SystemConfigKeyAcidNote                 = "acid_note"
	SystemConfigKeyTbtToken                 = "tbt_token"
	SystemConfigKeyBloodSugarNote           = "blood_sugar_note"
	SystemConfigKeyBmiNote                  = "bmi_note"
	SystemConfigKeyAppointmentTitle         = "appointment_title"
	SystemConfigKeyAppointmentDescription   = "appointment_description"
	SystemConfigKeyBloodFatNote             = "blood_fat_note"
	SystemConfigKeyRmyToken                 = "rmy_token"
	SystemConfigKeyMemberToken              = "member_token"
	SystemConfigKeyMemberBuyUrl             = "member_buy_url"
	SystemConfigKeyMemberQueryUrl           = "member_query_url"
	SystemConfigKeyChatCommandSwitch        = "chat_command_switch"
	SystemConfigKeyWecomGroupReplyMethod    = "wecom_group_reply_method"
	SystemConfigKeyWecomGroupReplyDelayTime = "wecom_group_reply_delay_time"
	SystemConfigKeyHezhenAppid              = "hezhen_appid"
	SystemConfigKeyHezhenSecret             = "hezhen_secret"
	SystemConfigKeyRmyPrivateKey            = "rmy_private_key"
	SystemConfigKeyRmyAccount               = "rmy_account"
	SystemConfigKeyRmyDomain                = "rmy_domain"
	SystemConfigKeyAgentProactiveChat       = "ai_agent_proactive_chat"
	SystemConfigKeyAgentAudioGenerator      = "ai_agent_audio_generator"
	SystemConfigKeyMissionRefreshTime       = "mission_refresh_time"
	SystemConfigKeyWorkflowCrawlArticle     = "ai_workflow_craw_article"
	SystemConfigKeyWorkflowGenWordsAndTips  = "ai_workflow_gen_words_tips"
	SystemConfigKeyLbsKeys                  = "lbs_keys"
	SystemConfigKeyTrackingPointKeys        = "tracking_point_keys"
	SystemConfigKeyArticleShow              = "article_show"
	SystemConfigKeyWorkflowAddFriend        = "ai_workflow_add_friend"
	SystemConfigKeyAppVersionLatest         = "app_version"
	SystemConfigKeyMeiSubcompanyList        = "mei_subcompany_list"
	SystemConfigKeyGlucoseTargetRange       = "glucose_target_range"
	SystemConfigKeyGroupSaleHour            = "group_sale_hour"
	SystemConfigKeyWorkflowGroupSale        = "ai_workflow_group_sale"
	SystemConfigKeyAutoSaleTypeList         = "auto_sale_type_list"
	SystemConfigKeyWorkflowAutoSale         = "ai_workflow_auto_sale"
	SystemConfigKeyWorkflowReportAnalysis   = "ai_workflow_report_analysis"
	SystemConfigKeyWorkflowReportAdvise     = "ai_workflow_report_advise"
	SystemConfigKeyWorkflowExcelChatMessage = "ai_workflow_excel_chat_message"
	SystemConfigKeyWorkflowFeishuChatApp    = "feishu_chat_app"
	SystemConfigKeyWorkflowReportProAdvise  = "ai_workflow_report_product_advise"
	SystemConfigKeyWorkflowKfQrcode         = "kefu_qrcode"
	SystemConfigKeyMeiHealthManager         = "mei_health_manager"
	SystemConfigKeyAiVisionAnalysis         = "ai_workflow_vision_analysis"
	SystemConfigKeyAiVisionImgConfig        = "ai_workflow_vision_img_config"
	SystemConfigKeyHealthGroupTypes         = "health_group_types"
	SystemConfigKeyTencentAsrConfig         = "tencent_asr_config"
	SystemConfigKeyAiFileDetect             = "ai_workflow_file_detect"      // 文件检测
	SystemConfigKeyPublicConfigKeys         = "public_config_keys"           // 公共配置项
	SystemConfigKeyWorkflowArticleAnalysis  = "ai_workflow_article_analysis" // 文章标签分析

	SystemConfigKeyWorkflowsemanticProcessing = "ai_workflow_semantic_processing" // 语义处理

	SystemConfigKeyWorkflowQuestAboutProblem = "ai_workflow_quest_about_problem" // 问卷相关问题AI
	SystemConfigKeyWorkflowQuestAboutAdvise  = "ai_workflow_quest_about_advose"  // 问卷相关建议AI
	SystemConfigKeyWorkflowMagicBox          = "ai_workflow_magic_box"           // 百宝箱 AI
	SystemConfigKeyWorkflowGeneral           = "ai_workflow_general"             // 通用AI

	SystemConfigKeyChronicCondition   = "chronic_condition"
	SystemConfigKeyHealthCheckInTypes = "health_check_in_types"
	SystemConfigKeyDefaultHabits      = "default_habits"

	SystemConfigKeyAppCancelAgreement  = "app_cancel_agreement"  // 注销协议
	SystemConfigKeyAppAboutUs          = "app_about_us"          // 关于我们
	SystemConfigKeyAppVersion          = "ahb_app_version"       // 版本号
	SystemConfigKeyHealthBeanRules     = "health_bean_rules"     // 健康豆规则
	SystemConfigKeyUserArchiveWorkflow = "user_archive_workflow" // 健康豆规则

	SystemConfigKeyHealthReportDaily         = "health_report_daily"         // 健康报告每日生成
	SystemConfigKeyLastweekHealthReport      = "lastweek_health_report"      // 健康报告每日生成
	SystemConfigKeyHealthRecordWorkflow      = "health_record_workflow"      // 图片膳食解析
	SystemConfigKeyHealthCategory            = "health_category"             // 健康信息采集分类
	SystemConfigKeyHealthScore               = "health_score"                // 健康信息采集分类
	SystemConfigKeyAgeHealthGuidelines       = "age_health_guidelines"       // 年龄健康指南
	SystemConfigKeyIndexBannerConfig         = "index_banner_config"         // 首页banner配置
	SystemConfigKeyPerfectDataConfig         = "perfect_data_config"         // 信息采集页面配置
	SystemConfigKeyMissionBaseConfig         = "mission_base_config"         // 每日任务配置
	SystemConfigKeyEihDefaultConfig          = "eih_default_config"          // 爱护网默认配置
	SystemConfigKeyEihSmsConfig              = "default_sms_config"          // 爱护网短信默认配置
	SystemConfigKeyDailyCheckinConfig        = "daily_checkin_config"        // 每日打卡任务配置
	SystemConfigKeyContentInteractionsConfig = "content_interactions_config" // 内容互动快捷短语配置
	SystemConfigKeyGlobalMessageConfig       = "global_message_config"       // 内容互动快捷短语配置
	SystemConfigKeySvcAdvisorConfig          = "svc_advisor_config"          // 服务顾问配置

	SystemConfigKeySkinTestAdditions = "skin_test_additions" // 测肤报告附加物

	IndexBannerHealthScore                    = "health_score"
	IndexBannerTotalRank                      = "total_rank"
	IndexBannerMcaAgedQs                      = "mca_aged_qs"
	SystemConfigKeyDynamicQuestionnaireConfig = "dynamic_questionnaire_config"

	GlobalMessageFieldHealthMetricsAlert = "health_metrics_alert"
	GlobalMessageFieldHabitReminder      = "habit_reminder"

	GlobalMessageFieldGroupJoinSuccess  = "group_join_success"
	GlobalMessageFieldGroupMemberJoined = "group_member_joined"
	GlobalMessageFieldGroupJoinRequest  = "group_join_request"
	GlobalMessageFieldGroupJoinApplied  = "group_join_applied"
	GlobalMessageFieldGroupJoinResult   = "group_join_result"
	GlobalMessageFieldGroupQuitResult   = "group_quit_result"

	SystemConfigKeyMpNoticeReortConf        = "mp_notice_report_config"        // 小程序通知配置 - 报告提醒
	SystemConfigKeyMpNoticeQuestReortConf   = "mp_notice_quest_report_config"  // 小程序通知配置 - 问卷报告提醒
	SystemConfigKeyMpNoticePhysicalExamConf = "mp_notice_physical_exam_config" // 小程序通知配置 - 体检提醒
	SystemConfigKeyMpNoticeUsagePlanConf    = "mp_notice_usage_plan_config"    // 小程序通知配置 - 产品计划提醒
	SystemConfigKeyMpNoticeProductUseConf   = "mp_notice_product_use_config"   // 小程序通知配置 - 产品使用提醒
	SystemConfigKeyMpNoticeClockInConf      = "mp_notice_clock_in_config"      // 小程序通知配置 - 打卡提醒

	SystemConfigKeyWxOfficialBloodPressureConf = "wx_official_blood_pressure_conf" // 公众号配置 - 血压
	SystemConfigKeyWxOfficialAcidConf          = "wx_official_acid_conf"           // 公众号配置 - 尿酸
	SystemConfigKeyWxOfficialBloodSugarConf    = "wx_official_blood_sugar_conf"    // 公众号配置 - 血糖
	SystemConfigKeyWxOfficialBloodFatConf      = "wx_official_blood_fat_conf"      // 公众号配置 - 血脂
	SystemConfigKeyWxOfficialBmiConf           = "wx_official_bmi_conf"            // 公众号配置 - BMI
	SystemConfigKeyWxOfficialQuestReortConf    = "wx_official_quest_report_conf"   // 公众号配置 - 问卷报告

	SystemConfigKeyMiniProgramConf = "miniprogram_conf" // 小程序配置
	SystemConfigKeyWxOfficialConf  = "wx_official_conf" // 公众号配置

	SystemConfigKeyWorkflowReportDynamic              = "ai_workflow_dynamic_report"
	SystemConfigKeyQuestionnaireImageUrls             = "questionnaire_image_urls"
	SystemConfigKeyDynamicQuestionnaireWorkflowConfig = "dynamic_questionnaire_workflow_config"
	SystemConfigKeyArticleSectionConfig               = "article_section_config"
	SystemConfigKeyNextSuggestedWorkflowConfig        = "next_suggested_workflow_config"
	SystemConfigKeyQuestionnaireStaticReportConfig    = "service_provider_report_general"
)

type (
	// Ai agent api访问设置
	SysConfAgentApi struct {
		ServerHost string `json:"server_host"`
		ApiKey     string `json:"api_key"`
	}

	// 小程序通知配置
	SysConfMpNotice struct {
		TemplateId string   `json:"template_id"` // 模板ID, 为空则默认不发送通知
		PagePath   string   `json:"page_path"`
		ParamKeys  []string `json:"param_keys"` // 参数key, 参数按显示顺序排列。输入参数将字段按顺序填入
	}

	// 服务商顾问配置
	SysConfSvcAdvisor struct {
		BindNoHint bool `json:"bind_no_hint"` // 无需弹出绑定提示
	}
)

var _ SystemConfigModel = (*customSystemConfigModel)(nil)

type (
	// SystemConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSystemConfigModel.
	SystemConfigModel interface {
		systemConfigModel
		FindByGroup(ctx context.Context, group string) ([]*SystemConfig, error)
		FindByKey(ctx context.Context, key string) (*SystemConfig, error)
	}

	customSystemConfigModel struct {
		*defaultSystemConfigModel
	}

	SkinTestAddition struct {
		Title string `json:"title"`
		Type  string `json:"type"`
		Value string `json:"value"`
	}
)

// 按传入类型获取配置
func GetConfigByKey[T any](ctx context.Context, sysModel SystemConfigModel, key string) (*T, error) {
	// 从数据库中获取配置
	systemConfig, err := sysModel.FindByKey(ctx, key)
	if err != nil {
		// 处理错误
		return nil, err
	}
	// 判断T的类型是否为string
	if _, ok := any(new(T)).(*string); ok {
		// 如果是string类型，直接返回Value字段
		return any(&systemConfig.Value).(*T), nil
	}
	// 解析配置
	var config T
	err = json.Unmarshal([]byte(systemConfig.Value), &config)
	if err != nil {
		// 处理错误
		return nil, err
	}
	return &config, nil
}

func (sc SystemConfig) GetValue() any {
	switch sc.ValueType {
	case "bool":
		parseVal, err := strconv.ParseBool(sc.Value)
		if err != nil {
			return nil
		}
		return parseVal
	case "int":
		parseVal, err := strconv.ParseInt(sc.Value, 10, 64)
		if err != nil {
			return nil
		}
		return parseVal
	case "list":
		var slice []string
		err := json.Unmarshal([]byte(sc.Value), &slice)
		if err != nil {
			return nil
		}
		return slice
	default:
		return sc.Value
	}
}

func (sc SystemConfig) GetHealthCheckInTypes() *HealthCheckInTypes {
	if sc.Key != SystemConfigKeyHealthCheckInTypes {
		return nil
	}
	var healthCheckInTypes HealthCheckInTypes
	err := json.Unmarshal([]byte(sc.Value), &healthCheckInTypes)
	if err != nil {
		return nil
	}
	return &healthCheckInTypes
}

func (sc SystemConfig) GetIndexBannerConfig() []*BannerConfig {
	if sc.Key != SystemConfigKeyIndexBannerConfig {
		return nil
	}
	var configs []*BannerConfig
	err := json.Unmarshal([]byte(sc.Value), &configs)
	if err != nil {
		return nil
	}
	return configs
}

func (sc SystemConfig) GetIndexBannerInfo(bannerName string) *BannerConfig {
	banners := sc.GetIndexBannerConfig()
	if banners == nil {
		return nil
	}
	for _, banner := range banners {
		if banner.Name == bannerName {
			return banner
		}
	}
	return nil
}

func (sc SystemConfig) GetEihUserHealthConfig() *EihDefaultConfig {
	if sc.Key != SystemConfigKeyEihDefaultConfig {
		return nil
	}
	var configs *EihDefaultConfig
	err := json.Unmarshal([]byte(sc.Value), &configs)
	if err != nil {
		return nil
	}
	return configs
}

func (sc SystemConfig) GetSmsConfig() *DefaultSmsConfig {
	if sc.Key != SystemConfigKeyEihSmsConfig {
		return nil
	}
	var configs *DefaultSmsConfig
	err := json.Unmarshal([]byte(sc.Value), &configs)
	if err != nil {
		return nil
	}
	return configs
}

func (sc SystemConfig) GetConfigByAge(age int) *AgeHealthConfig {
	if sc.Key != SystemConfigKeyAgeHealthGuidelines {
		return nil
	}

	// 将value字段解析为配置数组
	var configs []AgeHealthConfig
	if err := json.Unmarshal([]byte(sc.Value), &configs); err != nil {
		logx.Errorf("解析健康配置失败: %v", err)
		return nil
	}

	// 遍历查找匹配的年龄段
	for _, config := range configs {
		if age >= config.MinAge && age <= config.MaxAge {
			return &config
		}
	}
	return nil
}

func (sc SystemConfig) GetDynamicQuestionnaireConfig() *DynamicQuestionnaire {
	if sc.Key != SystemConfigKeyDynamicQuestionnaireConfig {
		return nil
	}
	var configs *DynamicQuestionnaire
	err := json.Unmarshal([]byte(sc.Value), &configs)
	if err != nil {
		return nil
	}
	return configs
}

func (sc SystemConfig) GetGlobalMessageConfigConfig() *GlobalMessageConfig {
	if sc.Key != SystemConfigKeyGlobalMessageConfig {
		return nil
	}
	var configs *GlobalMessageConfig
	err := json.Unmarshal([]byte(sc.Value), &configs)
	if err != nil {
		return nil
	}
	return configs
}

// GetGlobalMessageConfigValue 根据字段名获取全局消息配置中的值
func (sc SystemConfig) GetGlobalMessageConfigValue(name string) uint64 {
	config := sc.GetGlobalMessageConfigConfig()
	if config == nil {
		return 0
	}

	switch name {
	case GlobalMessageFieldHealthMetricsAlert:
		return config.HealthMetricsAlert
	case GlobalMessageFieldHabitReminder:
		return config.HabitReminder

	case GlobalMessageFieldGroupJoinSuccess:
		return config.HealthGroup.GroupJoinSuccess
	case GlobalMessageFieldGroupMemberJoined:
		return config.HealthGroup.GroupMemberJoined
	case GlobalMessageFieldGroupJoinRequest:
		return config.HealthGroup.GroupJoinRequest
	case GlobalMessageFieldGroupJoinApplied:
		return config.HealthGroup.GroupJoinApplied
	case GlobalMessageFieldGroupJoinResult:
		return config.HealthGroup.GroupJoinResult
	case GlobalMessageFieldGroupQuitResult:
		return config.HealthGroup.GroupQuitResult
	}

	return 0
}

// GetQuestionnaireImageUrls 获取问卷图片URL配置
func (sc SystemConfig) GetQuestionnaireImageUrls() *QuestionnaireImageUrls {
	if sc.Key != SystemConfigKeyQuestionnaireImageUrls {
		return nil
	}
	var config QuestionnaireImageUrls
	err := json.Unmarshal([]byte(sc.Value), &config)
	if err != nil {
		return nil
	}
	return &config
}

// NewSystemConfigModel returns a model for the database table.
func NewSystemConfigModel(conn *gorm.DB) SystemConfigModel {
	return &customSystemConfigModel{
		defaultSystemConfigModel: newSystemConfigModel(conn),
	}
}

func (m *customSystemConfigModel) FindByGroup(ctx context.Context, group string) ([]*SystemConfig, error) {
	var resp []*SystemConfig
	err := m.conn.WithContext(ctx).Model(&SystemConfig{}).Where("`group` = ?", group).Order("`key` ASC").Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customSystemConfigModel) FindByKey(ctx context.Context, key string) (*SystemConfig, error) {
	var resp SystemConfig
	err := m.conn.WithContext(ctx).Model(&SystemConfig{}).Where("`key` = ?", key).First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AiAgentConfig struct {
	ServerHost string `json:"server_host"`
	ApiKey     string `json:"api_key"`
}

type HealthTypes struct {
	ClassId     uint64 `json:"class_id"`
	Name        string `json:"name"`
	CheckInType uint64 `json:"check_in_type"`
	Icon        string `json:"icon"`
	CheckedIcon string `json:"checked_icon"`
	Description string `json:"description"`
}

type HealthCheckInTypes struct {
	Types []HealthTypes `json:"types"`
}

type BannerConfig struct {
	Name     string `json:"name"`
	JumpUrl  string `json:"jump_url"`
	JumpType int    `json:"jump_type"`
	ImgUrl   string `json:"img_url"`
}

type HealthScoreConfig struct {
	Comment string `json:"comment"`
	Stars   int    `json:"stars"`
	Level   string `json:"level"`
	Icon    string `json:"icon"`
}

type EihDefaultConfig struct {
	Icon            string `json:"icon"`
	Comment         string `json:"comment"`
	Level           string `json:"level"`
	DefaultNickname string `json:"default_nickname"`
	HasHabit        int    `json:"has_habit"`
	// 血压范围
	BloodPressure struct {
		SystolicMin  float64 `json:"systolic_min"`  // 收缩压最小值
		SystolicMax  float64 `json:"systolic_max"`  // 收缩压最大值
		DiastolicMin float64 `json:"diastolic_min"` // 舒张压最小值
		DiastolicMax float64 `json:"diastolic_max"` // 舒张压最大值
	} `json:"blood_pressure"`
	// 血糖范围
	BloodGlucose struct {
		FastingMax        float64 `json:"fasting_max"`          // 空腹血糖最大值
		PostPrandialMax   float64 `json:"post_prandial_max"`    // 餐后血糖最大值
		DiabetesLowMin    float64 `json:"diabetes_low_min"`     // 糖尿病患者低血糖最小值
		NonDiabetesLowMin float64 `json:"non_diabetes_low_min"` // 非糖尿病患者低血糖最小值
	} `json:"blood_glucose"`
	// 尿酸范围
	UricAcid struct {
		MaleMax   float64 `json:"male_max"`   // 男性尿酸最大值
		FemaleMax float64 `json:"female_max"` // 女性尿酸最大值
		Min       float64 `json:"min"`        // 尿酸最小值
		HMin      float64 `json:"h_min"`      // 痛风患者尿酸最小值
		HMax      float64 `json:"h_max"`      // 痛风患者尿酸最大值
	} `json:"uric_acid"`
	// 血脂范围
	BloodLipids struct {
		Max  float64 `json:"max"`   // 总胆固醇最大值
		Min  float64 `json:"min"`   // 总胆固醇最小值
		HMax float64 `json:"h_max"` // 高血脂患者总胆固醇最大值
	} `json:"blood_lipids"`
	// 各项指标扣分率
	PenaltyRate struct {
		BloodPressure float64 `json:"blood_pressure"` // 血压扣分率
		BloodGlucose  float64 `json:"blood_glucose"`  // 血糖扣分率
		UricAcid      float64 `json:"uric_acid"`      // 尿酸扣分率
		BloodLipids   float64 `json:"blood_lipids"`   // 血脂扣分率
	} `json:"penalty_rate"`
	// BMI相关参数
	BMI struct {
		Min         float64 `json:"min"`          // BMI最小值
		Max         float64 `json:"max"`          // BMI最大值
		PenaltyRate float64 `json:"penalty_rate"` // BMI扣分率
		MinScore    float64 `json:"min_score"`    // BMI最低分数
	} `json:"bmi"`
}

type AgeHealthConfig struct {
	MinAge     int `json:"min_age"`
	MaxAge     int `json:"max_age"`
	SleepMin   int `json:"sleep_min"`
	SleepMax   int `json:"sleep_max"`
	CalorieMin int `json:"calorie_min"`
	CalorieMax int `json:"calorie_max"`
}

type MissionBaseConfig struct {
	Title         string `json:"title"`
	Label         string `json:"label"`
	JumpType      int    `json:"jump_type"`
	JumpUrl       string `json:"jump_url"`
	Icon          string `json:"icon"`
	CompletedIcon string `json:"completed_icon"`
}

type DefaultSmsConfig struct {
	DefaultConfig string       `json:"default_config"`
	Configs       []SmsConfigs `json:"configs"`
}

type SmsConfigs struct {
	Name   string          `json:"name"`
	Config SmsConfigDetail `json:"config"`
}

type SmsConfigDetail struct {
	AccesskeyId     string           `json:"accesskey_id,omitempty"`
	AccesskeySecret string           `json:"accesskey_secret,omitempty"`
	Region          string           `json:"region"`
	BizType         []DefaultBizType `json:"biz_type"`
	SecretId        string           `json:"secret_id,omitempty"`
	SecretKey       string           `json:"secret_key,omitempty"`
	AppId           string           `json:"app_id,omitempty"`
}

type DefaultBizType struct {
	Name       []string `json:"name"`
	SignName   string   `json:"sign_name"`
	TemplateId string   `json:"template_id"`
}

type DailyTask struct {
	ID            uint64 `json:"id"`
	Label         string `json:"label"`
	DailyLimit    int64  `json:"daily_limit"`
	RequiredCount int    `json:"required_count"`
}

type CompletionRules struct {
	AllTasksRequired bool  `json:"all_tasks_required"`
	BonusPoints      int64 `json:"bonus_points"`
}

type DailyCheckinConfig struct {
	Tasks           []DailyTask      `json:"tasks"`
	CompletionRules *CompletionRules `json:"completion_rules"`
	MissionBase     *MissionBase     `json:"mission_base"`
}

type MissionBase struct {
	Id    uint64 `json:"id"`
	Label string `json:"label"`
}

type DynamicQuestionnaire struct {
	DefaultQuestionnaireTitle string `json:"default_questionnaire_title"`
	Modules                   struct {
		Immune        string `json:"immune"`
		AntiAging     string `json:"anti_aging"`
		BloodQi       string `json:"blood_qi"`
		SkinWhitening string `json:"skin_whitening"`
		LiverKidney   string `json:"liver_kidney"`
	} `json:"modules"`
	DefaultQuestionnaireId int `json:"default_questionnaire_id"`
	MinModuleItemCount     int `json:"min_module_item_count"`
	MaxModuleItemCount     int `json:"max_module_item_count"`
}

type GlobalMessageConfig struct {
	HealthGroup        HealthGroupConfig `json:"health_group"`
	HealthMetricsAlert uint64            `json:"health_metrics_alert"`
	HabitReminder      uint64            `json:"habit_reminder"`
}

type HealthGroupConfig struct {
	GroupJoinSuccess  uint64 `json:"group_join_success"`
	GroupMemberJoined uint64 `json:"group_member_joined"`
	GroupJoinRequest  uint64 `json:"group_join_request"`
	GroupJoinApplied  uint64 `json:"group_join_applied"`
	GroupJoinResult   uint64 `json:"group_join_result"`
	GroupQuitResult   uint64 `json:"group_quit_result"`
}

type QuestionnaireImageUrls struct {
	LiverCleanse        string `json:"liver_cleanse"`          // 肝阳
	NadPlus             string `json:"nad_plus"`               // NAD+
	Metabolism          string `json:"metabolism"`             // 代谢
	EnergyBoost         string `json:"energy_boost"`           // 能量
	BrighteningAntiage  string `json:"brightening_antiage"`    // 美白
	Doctor              string `json:"doctor"`                 // 医生
	DoctorBackground    string `json:"doctor_background"`      // 医生
	HealthGoals         string `json:"health_goals"`           // 健康目标
	HealthGoalsSubImage string `json:"health_goals_sub_image"` // 健康目标子图
	DifferenceIcon      string `json:"difference_icon"`        // 差异图标
}

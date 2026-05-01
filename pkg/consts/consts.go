package consts

const (
	FlagFalse uint64 = 0
	FlagTrue  uint64 = 1

	FlagEnableStr  string = "enable"
	FlagDisableStr string = "disable"

	StatusDaft       int64 = -1
	StatusNew        int64 = 0
	StatusProcessing int64 = 1
	StatusSuccessful int64 = 2
	StatusFailure    int64 = 3

	AccessToken string = "%s:%s:access:token:%s"

	DailyTipGenCountKey   = "tips:newcount"
	DailyWordsGenCountKey = "words:newcount"

	HealthRank  = "%s:%s:health:rank"
	InvitesRank = "%s:%s:invites:rank"

	AppNameDefault        = "" // 默认app name
	UserProfileKeyArchive = "user:profile:health:archive"

	ContentTemplateHomePage = "home_page"

	FeishuSheetTypeMaomei = 1
	FeishuSheetTypeMei    = 2

	DiabetesCategory      = "diabetes-special"
	BloodPressureCategory = "blood-pressure-survey"
	DiabetesSurvey        = "diabetes-survey"

	HealthEventEditModeAll     = "all"     // 所有事件
	HealthEventEditModeAfter   = "after"   // 当前及之后事件
	HealthEventEditModeCurrent = "current" // 当前事件

	HealthEventRedisKey = "health:event:uuid"
)

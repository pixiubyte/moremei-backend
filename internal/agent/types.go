package agent

type (
	DifyCommonOutput struct {
		Output string `json:"output"` // 使用 output 作为通用返回
		Error  string `json:"error"`  // 使用 error 作为通用错误返回
	}

	DifyArticleTagsAnalyzeResult struct {
		HealthTags  []string              `json:"health_tags"`
		SectionTags []map[string][]string `json:"section_tags"` // [ {"0": ["食品安全", "肠道健康",]} ]
	}
)

import (
	"context"
	"database/sql"
	"strings"
	{{if .time}}"time"{{end}}

	"yqunmeta/ai-saas/common/gormc"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stringx"
	"gorm.io/gorm"
)

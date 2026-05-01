package {{.pkg}}
{{if .withCache}}
import (
	"moremmei/ai-saas/common/gormc"
	"moremei/ai-saas/internal/multitenancy"

	"gorm.io/gorm"
)
{{else}}
import (
    "moremei/ai-saas/internal/multitenancy"

	"gorm.io/gorm"
)
{{end}}
var _ {{.upperStartCamelObject}}Model = (*custom{{.upperStartCamelObject}}Model)(nil)

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
		SwitchTenant(tenant *multitenancy.Tenant) {{.upperStartCamelObject}}Model
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(conn *gorm.DB{{if .withCache}}, c *gormc.CacheConf{{end}}) {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(conn{{if .withCache}}, c{{end}}),
	}
}

func (m *custom{{.upperStartCamelObject}}Model) SwitchTenant(tenant *multitenancy.Tenant) {{.upperStartCamelObject}}Model {
	return New{{.upperStartCamelObject}}Model(tenant.GetDB())
}

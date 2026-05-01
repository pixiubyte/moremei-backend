
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}, tx *gorm.DB) error {
    var err error
	{{if .withCache}}{{.keys}}
	if tx != nil {
        err = m.ExecTransCtx(ctx, tx, func(conn *gorm.DB) error {
            return conn.Create(&data).Error
        }, {{.keyValues}})
	} else {
        err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
            return conn.Create(&data).Error
        }, {{.keyValues}})
    }{{else}}if tx != nil {
        err = tx.WithContext(ctx).Create(&data).Error
    } else {
        err = m.conn.WithContext(ctx).Create(&data).Error
    }{{end}}
	return err
}
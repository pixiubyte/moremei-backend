
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, data *{{.upperStartCamelObject}}, tx *gorm.DB) error {
    var err error
	{{if .withCache}}{{.keys}}
	if tx != nil {
	    err = m.ExecTransCtx(ctx, tx, func(conn *gorm.DB) error {
    		return conn.Save(data).Error
    	}, {{.keyValues}})
	} else {
	    err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
    		return conn.Save(data).Error
    	}, {{.keyValues}})
	}{{else}}if tx != nil {
	    err = tx.WithContext(ctx).Save(data).Error
	} else {
	    err = m.conn.WithContext(ctx).Save(data).Error
	}{{end}}
	return err
}

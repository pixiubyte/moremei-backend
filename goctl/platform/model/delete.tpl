func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}, tx *gorm.DB) error {
    var err error
	{{if .withCache}}{{if .containsIndexCache}}data, err := m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    if tx != nil {
        err = m.ExecTransCtx(ctx, tx, func(conn *gorm.DB) error {
    	    return conn.Delete(&{{.upperStartCamelObject}}{}, {{.lowerStartCamelPrimaryKey}}).Error
    	}, {{.keyValues}})
    } else {
        err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
    	    return conn.Delete(&{{.upperStartCamelObject}}{}, {{.lowerStartCamelPrimaryKey}}).Error
    	}, {{.keyValues}})
    }{{else}}if tx != nil {
        err = tx.WithContext(ctx).Delete(&{{.upperStartCamelObject}}{}, {{.lowerStartCamelPrimaryKey}}).Error
    } else {
        err = m.conn.WithContext(ctx).Delete(&{{.upperStartCamelObject}}{}, {{.lowerStartCamelPrimaryKey}}).Error
    }{{end}}

	return err
}
package gormc

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

func NewClient(config *Config) (client *gorm.DB, error error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return
	}

	config.FillWithDefaults()

	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		PrepareStmt: config.PrepareStmt,
		Logger: NewLogger(logx.WithContext(context.Background()), gl.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  gl.Warn,
		}),
	})
	if err != nil {
		return
	}
	rawDB, err := db.DB()
	if err != nil {
		return
	}
	if err = rawDB.Ping(); err != nil {
		return
	}

	if config.MaxOpenConns > 0 {
		rawDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		rawDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxLifeConns > 0 {
		rawDB.SetConnMaxLifetime(time.Duration(config.MaxLifeConns) * time.Second)
	}

	if config.DebugSQL {
		db = db.Debug()
	}

	return db, nil
}

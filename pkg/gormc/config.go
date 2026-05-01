package gormc

import (
    "runtime"
    "time"

    "github.com/go-sql-driver/mysql"
)

// Config Config代表mysql实例的配置信息
type Config struct {
    Driver       string        `json:"Driver,optional"`
    DSN          string        `json:"DSN"`
    DialTimeout  time.Duration `json:",optional"`
    ReadTimeout  time.Duration `json:",optional"`
    WriteTimeout time.Duration `json:",optional"`
    MaxOpenConns int           `json:",optional"`
    MaxIdleConns int           `json:",optional"`
    MaxLifeConns int           `json:",optional"`
    DebugSQL     bool          `json:",optional"`
    PrepareStmt  bool          `json:",optional"`
}

// NewConfig returns a *mysql.Config with timeout settings
func (c *Config) NewConfig() (dsn *mysql.Config, err error) {
    dsn, err = mysql.ParseDSN(c.DSN)
    if err != nil {
        return
    }

    // adjust timeout of DSN
    if dsn.Timeout <= 0 {
        dsn.Timeout = c.DialTimeout * time.Millisecond
    }
    if dsn.ReadTimeout <= 0 {
        dsn.ReadTimeout = c.ReadTimeout * time.Millisecond
    }
    if dsn.WriteTimeout <= 0 {
        dsn.WriteTimeout = c.WriteTimeout * time.Millisecond
    }

    // sync
    c.DSN = dsn.FormatDSN()

    return
}

// FillWithDefaults apply default values for field with invalid value.
func (c *Config) FillWithDefaults() {
    maxCPU := runtime.NumCPU()

    if c.DialTimeout <= 0 || c.DialTimeout > time.Duration(MaxDialTimeout*maxCPU) {
        c.DialTimeout = MaxDialTimeout
    }

    if c.ReadTimeout <= 0 || c.ReadTimeout > time.Duration(MaxReadTimeout*maxCPU) {
        c.ReadTimeout = MaxReadTimeout
    }

    if c.WriteTimeout <= 0 || c.WriteTimeout > time.Duration(MaxWriteTimeout*maxCPU) {
        c.WriteTimeout = MaxWriteTimeout
    }

    if c.MaxOpenConns <= 0 || c.MaxOpenConns > MaxOpenConn*maxCPU {
        c.MaxOpenConns = MaxOpenConn
    }

    if c.MaxIdleConns <= 0 || c.MaxIdleConns > MaxIdleConn*maxCPU {
        c.MaxIdleConns = MaxIdleConn
    }

    if c.MaxLifeConns <= 0 || c.MaxLifeConns > MaxLifecycleConn*maxCPU {
        c.MaxLifeConns = MaxLifecycleConn
    }
}

type ManagerConfig map[string]*Config

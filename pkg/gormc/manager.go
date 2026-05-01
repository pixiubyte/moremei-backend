package gormc

import (
    "errors"
    "sync"

    "github.com/zeromicro/go-zero/core/logx"
    "gorm.io/gorm"
)

// Manager manages multi client instances of gorm for easy usage.
type Manager struct {
    clients sync.Map
    configs sync.Map
}

var (
    gormManager     *Manager
    gormManagerOnce sync.Once
)

func InitManager(config *ManagerConfig) {
    gormManagerOnce.Do(func() {
        gormManager = NewManger(config)
    })
}

func DefaultClient() *gorm.DB {
    return gormManager.NewClient(ClientMaster)
}

func Client(name string) *gorm.DB {
    return gormManager.NewClient(name)
}

// NewManger NewManager creates a new manager store of redis with configs.
func NewManger(configs *ManagerConfig) *Manager {
    mgr := &Manager{}
    mgr.Load(configs)

    return mgr
}

func (mgr *Manager) NewClient(name string) (client *gorm.DB) {
    mgrClient, ok := mgr.clients.Load(name)
    if ok {
        client, ok := mgrClient.(*gorm.DB)
        if ok {
            return client
        }
    }

    config, err := mgr.Config(name)
    if err != nil {
        logx.Errorf("get %s mysql config error %s", name, err.Error())
        return nil
    }

    client, err = NewClient(config)
    if err != nil {
        logx.Errorf("get %s mysql client error %s", name, err.Error())
        return nil
    }

    mgr.clients.Store(name, client)

    return client
}

// Config returns a config registered with the name given
func (mgr *Manager) Config(name string) (config *Config, err error) {
    mgrConfig, ok := mgr.configs.Load(name)
    if ok {
        config, ok := mgrConfig.(*Config)
        if ok {
            return config, nil
        }

        return nil, errors.New("no config found")
    }

    return nil, errors.New("named config is not a valid *Config type")
}

// Load registers all configs with its name defined by ManagerConfig
func (mgr *Manager) Load(configs *ManagerConfig) {
    if configs == nil {
        return
    }

    for name, config := range *configs {
        mgr.Add(name, config)
    }
}

func (mgr *Manager) Add(name string, config *Config) {
    config.FillWithDefaults()

    // store new config
    mgr.configs.Store(name, config)

    // remove old client
    mgr.clients.Delete(name)
}

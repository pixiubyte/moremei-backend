package cipher

import (
    "sync"

    "github.com/pkg/errors"
)

type Manager struct {
    ciphers sync.Map
    configs sync.Map
}

func NewManager(configs *ManagerConfig) *Manager {
    m := &Manager{}
    m.LoadConfigs(configs)

    return m
}

func (m *Manager) GetCipher(category string) (*Cipher, error) {
    return m.NewCipher(category)
}

func (m *Manager) NewCipher(category string) (*Cipher, error) {
    if cipher, ok := m.ciphers.Load(category); ok {
        if c, ok := cipher.(*Cipher); ok {
            return c, nil
        }
    }

    config, err := m.LoadConfig(category)
    if err != nil {
        return nil, err
    }

    cipher, err := NewCipher(category, config.Key, config.Iv)
    if err != nil {
        return nil, err
    }

    m.ciphers.Store(category, cipher)

    return cipher, nil
}

func (m *Manager) LoadConfig(category string) (*Config, error) {
    if config, ok := m.configs.Load(category); ok {
        if c, ok := config.(*Config); ok {
            return c, nil
        }

        return nil, errors.Errorf("Not a valid *Config type: %s", category)
    }

    return nil, errors.Errorf("No config found: %s", category)
}

func (m *Manager) LoadConfigs(configs *ManagerConfig) {
    if configs == nil {
        return
    }

    for name, config := range *configs {
        m.Add(name, config)
    }
}

func (m *Manager) Add(name string, config *Config) {
    // store new config
    m.configs.Store(name, config)

    // remove old client
    m.ciphers.Delete(name)
}

package config

import (
	"sync"
)

type GlobalConfig struct {
	config *Config
	BindAddress,
	LogLevel,
	ConfigPath,
	TrustRootDir string
	Mutex *sync.RWMutex
}

func NewGlobalConfig() *GlobalConfig {
	g := &GlobalConfig{}
	g.config = &Config{}
	g.Mutex = new(sync.RWMutex)
	return g
}

func (g *GlobalConfig) SetConfig(c *Config) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	g.config = c
}

func (g *GlobalConfig) GetConfig() *Config {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	return g.config
}

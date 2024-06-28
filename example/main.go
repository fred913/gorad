package main

import (
	"log/slog"
	"sync"

	"github.com/yiffyi/gorad"
	"github.com/yiffyi/gorad/data"
)

type Config struct {
	db          *data.JSONDatabase
	lock        *sync.RWMutex
	LogFileName string
	Debug       bool
}

func LoadConfig() *Config {
	lock := sync.RWMutex{}
	db := data.NewJSONDatabase("config.json", true)
	cfg := Config{
		db:   db,
		lock: &lock,
	}

	// will panic if file not exist
	db.Load(&cfg, true)
	return &cfg
}

func (c *Config) Save() {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.db.Save(c)

}

func main() {
	cfg := LoadConfig()

	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	logger := slog.New(gorad.NewTextFileSlogHandler(cfg.LogFileName, level))
	slog.SetDefault(logger)

	a := []int{1, 2, 3, 4, 5, 6}
	for k := range a {
		slog.Info("ranging", "k", k)
		a = a[0:1]
	}
}

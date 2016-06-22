package config

import (
	"time"
)

var (
	cfg     *DspConfig = nil
	cfgFile string
)

type ServerConfig struct {
	Listen  string
	MaxConn int `yaml:"max_conn,flow"`
}

type LogConfig struct {
	Path  string
	Level string
}

type RedisConfig struct {
	Addr    string `yaml:"addr,flow"`
	PoolNum int    `yaml:"pool_num,flow"`
}

type DspConfig struct {
	HttpServer ServerConfig `yaml:"http_server,flow"`
	Log        LogConfig    `yaml:"log,flow"`
	Redis      RedisConfig  `yaml:"redis,flow"`
	IpLib      string       `yaml:"ip_lib,flow"`
	Tag        TagConfig    `yaml:"tag"`
}

type TagConfig struct {
	Addr    string        `yaml:"addr,flow"`
	PoolNum int           `yaml:"pool_num,flow"`
	Timeout time.Duration `yaml:"timeout"`
}

func (p *DspConfig) GetLogPath() string {
	return p.Log.Path
}

func (p *DspConfig) GetLogLevel() string {
	return p.Log.Level
}

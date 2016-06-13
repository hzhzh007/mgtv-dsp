package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	//"time"
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
}

func (p *DspConfig) GetLogPath() string {
	return p.Log.Path
}

func (p *DspConfig) GetLogLevel() string {
	return p.Log.Level
}

func LoadConfig(configFile string) (config *DspConfig, err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	config = new(DspConfig)
	err = yaml.Unmarshal(data, config)
	return
}

func LoadAndSet(configFile string) (config *DspConfig, err error) {
	cfg, err = LoadConfig(configFile)
	if err != nil {
		return cfg, err
	}
	cfgFile = configFile
	return cfg, err
}

func GetConfig() *DspConfig {
	return cfg
}

func Reload() {
	c, err := LoadConfig(cfgFile)
	if err != nil {
		return
	}
	cfg = c
}

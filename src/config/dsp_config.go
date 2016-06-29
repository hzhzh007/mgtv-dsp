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

type MgtvConfig struct {
	BidPath         string `yaml:"bid_path"`
	NoticePath      string `yaml:"notice_path"`
	WinNoticeUrl    string `yaml:"win_notice_url"` //"http://wy.mgtv.com/winnotice?c=%%SETTLE_PRICE%%"
	ClickNonticeUrl string `yaml:"click_notice_url"`
	Key             string `yaml:"key"`
	RedirectHost    string `yaml:"redirect_host`
}

type DspConfig struct {
	HttpServer ServerConfig  `yaml:"http_server,flow"`
	Log        LogConfig     `yaml:"log,flow"`
	Redis      RedisConfig   `yaml:"redis,flow"`
	IpLib      string        `yaml:"ip_lib,flow"`
	Tag        TagConfig     `yaml:"tag"`
	Mgtv       MgtvConfig    `yaml:"mgtv"`
	Resource   ResouceConfig `yaml:"resource"`
}

type ResouceConfig struct {
	Activity struct {
		Location string
		Reload   time.Duration
	}
	Flow struct {
		Reload time.Duration
	}
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

func LoadAndSet(configFile string) (config *DspConfig, err error) {
	cfg := new(DspConfig)
	err = LoadConfig(configFile, cfg)
	if err != nil {
		return cfg, err
	}
	cfgFile = configFile
	return cfg, err
}

func GetConfig() *DspConfig {
	return cfg
}

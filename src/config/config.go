package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	//"time"
)

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

func Reload() {
	c := new(DspConfig)
	err := LoadConfig(cfgFile, c)
	if err != nil {
		return
	}
	cfg = c
}
func LoadConfig(configFile string, config interface{}) (err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, config)
	return
}

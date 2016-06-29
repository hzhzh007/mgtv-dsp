package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	//"time"
)

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
	return LoadConfigFromBytes(data, config)
}

func LoadConfigFromBytes(data []byte, config interface{}) error {
	return yaml.Unmarshal(data, config)
}

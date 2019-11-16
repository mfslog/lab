package config

import (
	"github.com/spf13/viper"
)

type config struct{
	Endpoint string `yaml:"endpoint"`
	RemotePath string `yaml:"remotePath"`
	ProviderType string `yaml:"providerType"`
}

func Init()(err error){
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yaml")
	cfg := &config{}
	err = viper.Unmarshal(cfg)
	if err != nil{
		return err
	}
	err = viper.AddRemoteProvider(cfg.ProviderType, cfg.Endpoint, cfg.RemotePath)
	if err != nil{
		return err
	}
	return nil
}
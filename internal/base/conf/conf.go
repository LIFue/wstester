package conf

import (
	"github.com/spf13/viper"
	"wstester/internal/base/data"
	"wstester/internal/base/server"
	"wstester/pkg/log"
)

const DefaultConfigFileName = "../../../config/config.yaml"

type AllConfig struct {
	Debug  bool    `json:"debug" mapstructure:"debug" yaml:"debug"`
	Data   *Data   `json:"data" mapstructure:"data" yaml:"data"`
	Server *Server `json:"server" mapstructure:"server" yaml:"server"`
}

type Data struct {
	Database *data.Database `json:"database" mapstructure:"database" yaml:"database"`
}

type Server struct {
	WebSocket *server.WebSocket `json:"web_socket" mapstructure:"web_socket" yaml:"web_socket"`
}

// ReadConfig read config
func ReadConfig(configFilePath string) (c *AllConfig, err error) {
	c = &AllConfig{}

	if len(configFilePath) == 0 {
		configFilePath = DefaultConfigFileName
	}
	log.Infof("config file path: %s", configFilePath)
	viper.SetConfigFile(configFilePath)
	if err = viper.ReadInConfig(); err != nil {
		log.Errorf("read config error: %s", err.Error())
		return nil, err
	}
	if err = viper.Unmarshal(c); err != nil {
		return nil, err
	}
	return c, nil
}

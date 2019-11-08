package config

import (
	"fmt"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"

	"github.com/jaylee630/Hobbit/utils/log"
)

type (
	BasicConfig struct {
		Host      string           `yaml:"host" json:"host"`
		Port      int              `yaml:"port" json:"port"`
		AdminPort int              `yaml:"admin_port" json:"admin_port"`
		MySQL     MySQLConfig      `yaml:"mysql" json:"mysql"`
		Logger    log.LoggerConfig `yaml:"logger" json:"logger"`
		Gateway   GatewayConfig    `yaml:"gateway" json:"gateway"`
	}

	GatewayConfig struct {
		Host string `yaml:"host" json:"host"`
		Port int    `yaml:"port" json:"port"`
	}

	MySQLConfig struct {
		Host        string `yaml:"host" json:"host"`
		Port        int    `yaml:"port" json:"port"`
		Username    string `yaml:"username" json:"username"`
		Password    string `yaml:"password" json:"password"`
		Database    string `yaml:"database" json:"database"`
		TablePrefix string `yaml:"table_prefix" json:"table_prefix" envconfig:"TABLE_PREFIX"`
	}

	Configurator interface {
		GetAppAddr() string
		GetGatewayAddr() string
		GetDBSource() map[string]string
	}
)

func (n *BasicConfig) GetAppAddr() string {
	return fmt.Sprintf("http://%s:%d", n.Host, n.Port)
}

func (n *BasicConfig) GetGatewayAddr() string {
	return fmt.Sprintf("http://%s:%d", n.Gateway.Host, n.Gateway.Port)
}

func (n *BasicConfig) GetDBSource() map[string]string {
	return map[string]string{
		"mysql": fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true",
			n.MySQL.Username, n.MySQL.Password,
			n.MySQL.Host, n.MySQL.Port, n.MySQL.Database),
	}
}

func loadYAMLConfig(data []byte, cfg *BasicConfig) error {

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config content, error: %s ", err.Error())
	}

	return nil
}

func loadEnvConfig(prefix string, cfg *BasicConfig) error {

	if err := envconfig.Process(prefix, cfg); err != nil {
		return err
	}

	if err := envconfig.Process(prefix+"_MYSQL", &cfg.MySQL); err != nil {
		return err
	}

	if err := envconfig.Process(prefix+"_GATEWAY", &cfg.Gateway); err != nil {
		return err
	}

	if err := envconfig.Process(prefix+"_LOGGER", &cfg.Logger); err != nil {
		return err
	}

	return nil
}

func LoadConfig(path, envPrefix string) (*BasicConfig, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &BasicConfig{}
	err = loadYAMLConfig(data, cfg)
	if err != nil {
		return nil, err
	}

	err = loadEnvConfig(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type App struct {
	Name        string `yaml:"name"`
	Debug       bool   `yaml:"debug"`
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
}

type Jwt struct {
	AccessTokenExpTime  time.Duration `yaml:"access_token_expired_time"`
	RefreshTokenExpTime time.Duration `yaml:"refresh_token_expired_time"`
}

type Metrics struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

// Jaeger - contains all parameters metrics information.
type Jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type MongoDB struct {
	Dsn        string `yaml:"dsn"`
	Address    string `yaml:"addr"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

type HTTP struct {
	Address     string `yaml:"address"`
	Port        string `yaml:"port"`
	SwaggerPort string `yaml:"swagger_port"`
	Profile     string `yaml:"profile"`
}

type Grpc struct {
	Address           string `yaml:"address"`
	Port              string `yaml:"port"`
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
	Timeout           time.Duration
}

type Configs struct {
	App     App
	Jwt     Jwt
	HTTP    HTTP
	Metrics Metrics
	MongoDB MongoDB
	Jaeger  Jaeger
	Grpc    Grpc
}

func NewConfigs() (*Configs, error) {
	cfg := &Configs{}
	if err := cfg.InitConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Configs) InitConfig() error {
	if err := c.GetConfig("config.yml"); err != nil {
		return err
	}
	return nil
}

func (c *Configs) GetHTTPAddr() string {
	return c.HTTP.Address + ":" + c.HTTP.Port
}

func (c *Configs) GetConfig(cfgName string) error {
	if err := c.loadFile("./" + cfgName); err != nil {
		return err
	}
	return nil
}

func (c *Configs) GetConnStr(dbName string) string {
	var res string
	switch dbName {
	case "postgres":
	case "mongo":
		res = " mongodb://" + os.Getenv("MONGO_LOGIN") + ":" + os.Getenv("MONGO_PASS") + "@" + c.MongoDB.Address + ":" + c.MongoDB.Port + "/"
	default:
		return ""
	}
	return res
}

func (c *Configs) loadFile(name string) error {
	file, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return err
	}
	return nil
}

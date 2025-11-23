package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env    string       `mapstructure:"env"`
	DB     DBConfig     `mapstructure:"db"`
	Server ServerConfig `mapstructure:"server"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DSN constructs the connection string automatically
// Format: user:password@tcp(host:port)/dbname?parseTime=true
func (db DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("./config")

	// 1. Enable Environment Variable Overrides
	viper.AutomaticEnv()

	// 2. IMPORTANT: Replace dots with underscores
	// This allows you to use DB_HOST env var to override db.host in yaml
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

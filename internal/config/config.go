package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type MainConfig struct {
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
}

type MysqlConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
}

type JWTConfig struct {
	Secret      string `toml:"secret"`
	Issuer      string `toml:"issuer"`
	Audience    string `toml:"audience"`
	ExpireHours int    `toml:"expireHours"`
}

type LogConfig struct {
	LogPath string `toml:"logPath"`
}

type KafkaConfig struct {
	MessageMode string        `toml:"messageMode"`
	HostPort    string        `toml:"hostPort"`
	LoginTopic  string        `toml:"loginTopic"`
	LogoutTopic string        `toml:"logoutTopic"`
	ChatTopic   string        `toml:"chatTopic"`
	Partition   int           `toml:"partition"`
	Timeout     time.Duration `toml:"timeout"`
}

type StaticSrcConfig struct {
	StaticAvatarPath string `toml:"staticAvatarPath"`
}

type Config struct {
	MainConfig      `toml:"mainConfig"`
	MysqlConfig     `toml:"mysqlConfig"`
	RedisConfig     `toml:"redisConfig"`
	JWTConfig       `toml:"jwtConfig"`
	LogConfig       `toml:"logConfig"`
	KafkaConfig     `toml:"kafkaConfig"`
	StaticSrcConfig `toml:"staticSrcConfig"`
}

var config *Config

func LoadConfig() error {
	configPath := os.Getenv("KAMA_CHAT_CONFIG")
	if configPath == "" {
		configPath = "configs/config.toml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			_, sourceFile, _, ok := runtime.Caller(0)
			if ok {
				configPath = filepath.Join(filepath.Dir(sourceFile), "..", "..", "configs", "config.toml")
			}
		}
	}
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return fmt.Errorf("load config %q: %w", configPath, err)
	}
	if secret := os.Getenv("KAMA_CHAT_JWT_SECRET"); secret != "" {
		config.JWTConfig.Secret = secret
	}
	if config.JWTConfig.ExpireHours <= 0 {
		config.JWTConfig.ExpireHours = 24
	}
	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		if err := LoadConfig(); err != nil {
			panic(err)
		}
	}
	return config
}

package config

import (
	"manga-go/internal/pkg/logger"

	"github.com/spf13/viper"
)

type postgresqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type jwtConfig struct {
	Secret        string `mapstructure:"secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
	ExpiresAt     uint32 `mapstructure:"expires_seconds"`
	RefreshExpire uint32 `mapstructure:"refresh_expire_seconds"`
}

type OtlpConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type ServiceConfig struct {
	Name         string `mapstructure:"name"`
	Domain       string `mapstructure:"domain"`
	Port         int    `mapstructure:"port"`
	AllowOrigins string `mapstructure:"allow_origins"`
	DebugMode    bool   `mapstructure:"debug_mode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Cluster  string `mapstructure:"cluster"`
}

type smtpConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	FromName string `mapstructure:"from_name"`
	FromMail string `mapstructure:"from_mail"`
}

type asynqConfig struct {
	Concurrency int `mapstructure:"concurrency"`
}

type CookieNameConfig struct {
	AccessToken  string `mapstructure:"access_token"`
	RefreshToken string `mapstructure:"refresh_token"`
}

type ResetPasswordConfig struct {
	TokenExpiryMinutes int    `mapstructure:"token_expiry_minutes"`
	ResetPasswordURL   string `mapstructure:"reset_password_url"`
}

type ObjectStorageConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key"`
	SecretAccessKey string `mapstructure:"secret_key"`
	BucketName      string `mapstructure:"bucket_name"`
	Region          string `mapstructure:"region"`
}

type SeederConfig struct {
	AdminEmail    string `mapstructure:"admin_email"`
	AdminPassword string `mapstructure:"admin_password"`
}

type Config struct {
	Production    bool                `mapstructure:"production"`
	PostgreSQL    postgresqlConfig    `mapstructure:"db"`
	Jwt           jwtConfig           `mapstructure:"jwt"`
	Otlp          OtlpConfig          `mapstructure:"otlp"`
	Service       ServiceConfig       `mapstructure:"service"`
	Redis         RedisConfig         `mapstructure:"redis"`
	SMTP          smtpConfig          `mapstructure:"smtp"`
	Asynq         asynqConfig         `mapstructure:"asynq"`
	CookieName    CookieNameConfig    `mapstructure:"cookie_name"`
	ResetPassword ResetPasswordConfig `mapstructure:"reset_password"`
	ObjectStorage ObjectStorageConfig `mapstructure:"object_storage"`
	Seeder        SeederConfig        `mapstructure:"seeder"`
}

func LoadConfig() *Config {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile("config.yml")
	v.SetConfigType("yml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logger.GetLogger().Fatalf("Error while reading config file: %v", err)
		panic(err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		logger.GetLogger().Fatalf("Error while unmarshaling config file: %v", err)
		panic(err)
	}

	return &config
}

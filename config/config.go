package config

import (
	"github.com/spf13/viper"
	"gopkg.in/validator.v2"
)

type ConfigEnv struct {
	Stage                  string `mapstructure:"GO_AUTH_STAGE"                     validate:"nonzero"`
	Port                   string `mapstructure:"GO_AUTH_PORT"                      validate:"nonzero"`
	AdminSecret            string `mapstructure:"GO_AUTH_ADMIN_SECRET"              validate:"nonzero,len=32"`
	JwtTokenSecret         string `mapstructure:"GO_AUTH_JWT_TOKEN_SECRET"          validate:"nonzero,len=32"`
	JwtRefreshTokenSecret  string `mapstructure:"GO_AUTH_JWT_REFRESH_TOKEN_SECRET"  validate:"nonzero,len=32"`
	JwtTokenExpired        string `mapstructure:"GO_AUTH_JWT_TOKEN_EXPIRED"         validate:"nonzero"`
	JwtRefreshTokenExpired string `mapstructure:"GO_AUTH_JWT_REFRESH_TOKEN_EXPIRED" validate:"nonzero"`
	DataBaseHost           string `mapstructure:"GO_AUTH_DB_HOST"                   validate:"nonzero"`
	DataBasePort           string `mapstructure:"GO_AUTH_DB_PORT"                   validate:"nonzero"`
	DataBaseName           string `mapstructure:"GO_AUTH_DB_NAME"                   validate:"nonzero"`
	DataBaseUser           string `mapstructure:"GO_AUTH_DB_USER"                   validate:"nonzero"`
	DataBasePassword       string `mapstructure:"GO_AUTH_DB_PASS"                   validate:"nonzero"`
	DataBaseAutoMigrate    bool   `mapstructure:"GO_AUTH_DB_AUTO_MIGRATE"`
	TicketExpiresIn        string `mapstructure:"GO_AUTH_TICKET_EXPIRES_IN"`
}

func ConfigService() (configEnv ConfigEnv) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetDefault("GO_AUTH_STAGE", "develop")
	viper.SetDefault("GO_AUTH_PORT", "3000")
	viper.SetDefault("GO_AUTH_JWT_TOKEN_SECRET", "hello")
	viper.SetDefault("GO_AUTH_JWT_REFRESH_TOKEN_SECRET", "world")
	viper.SetDefault("GO_AUTH_JWT_TOKEN_EXPIRED", "30m")
	viper.SetDefault("GO_AUTH_JWT_REFRESH_TOKEN_EXPIRED", "24h")
	viper.SetDefault("GO_AUTH_DB_PORT", "5432")
	viper.SetDefault("GO_AUTH_DB_AUTO_MIGRATE", false)
	viper.SetDefault("GO_AUTH_TICKET_EXPIRES_IN", "1h")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&configEnv)
	if err != nil {
		panic(err)
	}

	err = validator.Validate(configEnv)
	if err != nil {
		panic(err)
	}

	return configEnv
}

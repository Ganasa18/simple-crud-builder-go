package appconfig

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/Ganasa18/simple-crud-builder-go/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppUrl                   string
	AppEnv                   string
	AppDebug                 string
	AppVersion               string
	AppName                  string
	AppPort                  string
	AppDir                   string
	GinMode                  string
	TemplateDir              string
	DbTZ                     string
	BasicAuthUsername        string
	BasicAuthPassword        string
	DbHost                   string
	DbPort                   int
	DbName                   string
	DbUsername               string
	DbPass                   string
	DbSchema                 string
	DbDriver                 string
	JWTSecretAccessToken     string
	AuthExpiredToken         string
	AuthExpiredTokenDuration time.Duration
}

func (c Config) IsStaging() bool {
	return c.AppEnv == "development"
}

func (c Config) IsProd() bool {
	return c.AppEnv == "production"
}

func (c Config) IsDebug() bool {
	return c.AppDebug == "True"
}

func InitAppConfig() *Config {
	c := Config{}

	err := godotenv.Load()
	utils.IsErrorDoPanic(err)

	// check load type environment
	if os.Getenv(utils.CONFIG_APP_ENV) == "development" {
		logrus.Infoln("load devel environment variable ")
	} else if os.Getenv(utils.CONFIG_APP_ENV) == "staging" {
		logrus.Infoln("load staging environment variable ")
	} else if os.Getenv(utils.CONFIG_APP_ENV) == "production" {
		logrus.Infoln("load production environment variable ")
	}

	c.AppUrl = os.Getenv(utils.CONFIG_APP_URL)
	c.AppEnv = os.Getenv(utils.CONFIG_APP_ENV)
	c.AppDebug = os.Getenv(utils.CONFIG_APP_DEBUG)
	c.AppVersion = os.Getenv(utils.CONFIG_APP_VERSION)
	c.AppName = os.Getenv(utils.CONFIG_APP_NAME)
	c.AppPort = os.Getenv(utils.CONFIG_APP_PORT)
	c.GinMode = os.Getenv(utils.CONFIG_GIN_MODE)

	// db conn
	c.DbHost = os.Getenv(utils.CONFIG_DB_HOST)
	c.DbPort, _ = strconv.Atoi(os.Getenv(utils.CONFIG_DB_PORT))
	c.DbName = os.Getenv(utils.CONFIG_DB_NAME)
	c.DbUsername = os.Getenv(utils.CONFIG_DB_USERNAME)
	c.DbPass = os.Getenv(utils.CONFIG_DB_PASSWORD)
	c.DbSchema = os.Getenv(utils.CONFIG_DB_SCHEMA)
	c.DbDriver = os.Getenv(utils.CONFIG_DB_DRIVER)

	// jwt config
	c.JWTSecretAccessToken = os.Getenv(utils.CONFIG_JWT_SECRET_ACCESS_TOKEN)
	c.AuthExpiredToken = os.Getenv(utils.CONFIG_AUTH_EXPIRE_ACCESS_TOKEN)

	_, b, _, _ := runtime.Caller(0)
	appDir := path.Join(path.Dir(b), "..")
	c.AppDir = appDir
	c.TemplateDir = appDir + "/internal/template"

	return &c
}

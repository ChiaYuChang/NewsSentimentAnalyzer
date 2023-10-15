package global

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var AppVar Option

func ReadConfig() error {
	if err := BindFlags(); err != nil {
		return fmt.Errorf("error while Bindflags: %w", err)
	}

	if err := BindEnv(); err != nil {
		return fmt.Errorf("error while BindEnv: %w", err)
	}

	if !viper.GetBool("USE_DOCKER_ENV") {
		if err := ReadAndExportEnvFile("./.env"); err != nil {
			return fmt.Errorf("error while ReadAndExportEnvFile: %w", err)
		}
		viper.Set("POSTGRES_DB_NAME", fmt.Sprintf(
			"%s_%s", viper.GetString("APP_NAME"), viper.GetString("APP_STATE")),
		)
		viper.Set("REDIS_DB_NAME", fmt.Sprintf(
			"%s_%s", viper.GetString("APP_NAME"), viper.GetString("APP_STATE")),
		)
	}

	if err := ReadAndSetPostgresPassword(
		viper.GetString("POSTGRES_PASSWORD_FILE")); err != nil {
		return fmt.Errorf("error while reading postgres password: %w", err)
	}

	if err := ReadOptionsFrom(viper.GetString("APP_CONFIG_FILE")); err != nil {
		return fmt.Errorf("error while reading option file: %w", err)
	}

	if err := viper.Unmarshal(&AppVar); err != nil {
		return fmt.Errorf("error while unmarshaling AppVar: %w", err)
	}

	if err := ReadTokenMakerSecret(viper.GetString("TOKEN_SECRET_FILE")); err != nil {
		return fmt.Errorf("error while reading tokenmaker secret: %w", err)
	}
	return nil
}

func BindFlags() error {
	pflag.IntP("port", "p", 8000, "http server port")
	pflag.StringP("host", "h", "127.0.0.1", "http server host")
	pflag.StringP("version", "v", "v1", "api version")
	pflag.StringP("config", "c", "./config/option.json", "path to the configuration file")
	pflag.StringP("state", "s", "development", "the state of app")
	pflag.BoolP("docker", "d", false, "using docker env")
	pflag.Lookup("docker").NoOptDefVal = "true"
	pflag.Parse()

	viper.RegisterAlias("USE_DOCKER_ENV", "docker")
	viper.RegisterAlias("APP_HOST", "host")
	viper.RegisterAlias("APP_PORT", "port")
	viper.RegisterAlias("APP_STATE", "state")
	viper.RegisterAlias("APP_CONFIG_FILE", "config")
	viper.RegisterAlias("APP_API_VERSION", "version")
	return viper.BindPFlags(pflag.CommandLine)
}

func BindEnv() error {
	setDefaultEnvVariable()
	for _, env := range []string{
		"POSTGRES_USERNAME",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_SSL_MODE",
		"POSTGRES_PASSWORD",
		"POSTGRES_PASSWORD_FILE",
		"POSTGRES_DB_NAME",
		"TOKEN_SECRET_FILE",
		"APP_NAME",
	} {
		if err := viper.BindEnv(env); err != nil {
			return err
		}
	}
	return nil
}

func ReadAndExportEnvFile(path string) error {
	envfile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	envVars := strings.Split(string(envfile), "\n")
	for _, ev := range envVars {
		ev = strings.Trim(ev, "")
		if len(ev) > 0 {
			nvPair := strings.Split(ev, "=")
			if err := os.Setenv(nvPair[0], nvPair[1]); err != nil {
				return err
			}
		}
	}
	return nil
}

func ReadOptionsFrom(path string) error {
	fmt.Println("Read config from:", path)

	configFilePath, configFile := filepath.Split(
		filepath.Clean(path))
	configFileName := strings.Split(configFile, ".")

	if err := ReadOptions(configFileName[0], configFileName[1], configFilePath); err != nil {
		return err
	}
	return nil
}

func ReadOptions(configName, configType string, configPath ...string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	for _, cf := range configPath {
		viper.AddConfigPath(cf)
	}
	setDefaultForOption()
	return viper.ReadInConfig()
}

func isFileExist(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	// should not be an directory
	return !info.IsDir()
}

func ReadTokenMakerSecret(path string) error {
	viper.SetDefault("TOKEN_SECRET", "SHOULD-NEVER-USED-IN-PRODUCTION")
	AppVar.Token.SetSecret([]byte("SHOULD-NEVER-USED-IN-PRODUCTION"))

	if path != "" && isFileExist(path) {
		if secret, err := os.ReadFile(path); err != nil {
			return err
		} else {
			viper.Set("TOKEN_SECRET", string(secret))
			AppVar.Token.SetSecret(secret)
			return nil
		}
	} else {
		if viper.GetString("APP_STATE") == "production" {
			return os.ErrNotExist
		}
		return nil
	}
}

func ReadAndSetPostgresPassword(path string) error {
	if path != "" && isFileExist(path) {
		fmt.Printf("reading %s ...\n", path)
		if secret, err := os.ReadFile(path); err != nil {
			return err
		} else {
			viper.Set("POSTGRES_PASSWORD", string(secret))
			return nil
		}
	} else {
		if viper.GetString("APP_STATE") == "production" {
			return os.ErrNotExist
		}
		return nil
	}
}

func ReadOptionsFromFile(configFile io.Reader, configType string) error {
	viper.SetConfigType(configType)
	setDefaultForOption()
	return viper.ReadConfig(configFile)
}

func WriteConfig(filename string) error {
	return viper.SafeWriteConfigAs(filename)
}

func setDefaultEnvVariable() {
	viper.SetDefault("POSTGRES_USERNAME", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_HOST", "127.0.0.1")
	viper.SetDefault("POSTGRES_PORT", "5434")
	viper.SetDefault("POSTGRES_SSL_MODE", "disable")
	viper.SetDefault("POSTGRES_DB_NAME", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD_FILE", "")
	viper.SetDefault("TOKEN_SECRET_FILE", "/run/secrets/TOKEN_SECRET")

	viper.SetDefault("REDIS_HOST", "127.0.0.1")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB_NAME", "redis")
	viper.SetDefault("REDIS_NETWORK", "tcp")
	viper.SetDefault("REDIS_MAX_RETRIES", 5)
	viper.SetDefault("REDIS_READ_TIMEOUT", 3*time.Second)
	viper.SetDefault("REDIS_WRITE_TIMEOUT", 3*time.Second)
	viper.SetDefault("REDIS_FIFO", true)
	viper.SetDefault("REDIS_POOLSIZE", 6)
}

func setDefaultForOption() {
	viper.SetDefault("Token.SignMethod.Algorithm", "HMAC")
	viper.SetDefault("Token.SignMethod.Size", 384)
	viper.SetDefault("Token.ExpireAfter", 72*time.Hour)
	viper.SetDefault("Token.ValidAfter", 0*time.Second)

	viper.SetDefault("Password.ASCIIOnly", true)
	viper.SetDefault("Password.MaxLength", 32)
	viper.SetDefault("Password.MinLength", 8)
	viper.SetDefault("Password.MinNumDigit", 1)
	viper.SetDefault("Password.MinNumUpper", 1)
	viper.SetDefault("Password.MinNumLower", 1)
	viper.SetDefault("Password.MinNumSpecial", 1)

	viper.SetDefault("App.StaticFile.Path", "views/static")
	viper.SetDefault("App.StaticFile.SubFolder.image", "/image")
	viper.SetDefault("App.StaticFile.SubFolder.js", "/js")
	viper.SetDefault("App.StaticFile.SubFolder.css", "/css")
	viper.SetDefault("App.RoutePattern.StaticPage", "/static/*")

	viper.SetDefault("App.Log", "./log.json")

	viper.SetDefault("App.SSL.Path", "./secrets")
	viper.SetDefault("App.SSL.CertFile", "server.crt")
	viper.SetDefault("App.SSL.KeyFile", ".server.key")
}

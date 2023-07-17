package global_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var jsonConfig = []byte(`{
  "token": {
    "signature": "jwt",
    "signMethod": {
      "size": 512
    },
    "expireAfter": "36h",
    "validAfter": "0s",
    "issuer": "",
    "subject": "",
    "audience": "",
    "id": ""
  },
  "password": {
    "minLength": 8,
    "maxLength": 30,
    "minNumDigit": 1,
    "minNumLower": 1,
    "minNumSpecial": 1
  },
  "app": {
    "staticFile": {
      "path": "views/static",
      "subfolder": {
        "image": "/image",
        "css": "/css"
      }
    },
    "template": [
      "views/template/*.gotmpl",
      "views/template/endpoint/*.gotmpl"
    ],
    "routePattern": {
      "page": {
        "login": "/login",
        "sign-up": "/sign-up",
        "sign-in": "/sign-in",
        "sign-out": "/sign-out",
        "admin": "/admin",
        "welcome": "/welcome",
        "apikey": "/apikey",
        "change-password": "/change-password",
        "endpoints": "/endpoints"
      },
      "errorPage": {
        "unauthorized": "/unauthorized",
        "bad-request": "/bad-request",
        "forbidden": "/forbidden"
      }
    }
  }
}`)

func TestReadConfig(t *testing.T) {
	var err error
	err = global.ReadOptionsFromFile(bytes.NewBuffer(jsonConfig), "json")
	require.NoError(t, err)

	var opt global.Option
	err = viper.Unmarshal(&opt)
	require.NoError(t, err)

	require.Equal(t, "HMAC", opt.Token.SignMethod.Algorthm)
	require.Equal(t, 36*time.Hour, opt.Token.ExpireAfter)
	require.Equal(t, 0*time.Hour, opt.Token.ValidAfter)
	require.Equal(t, 1, opt.Password.MinNumUpper)
	require.Equal(t, true, opt.Password.ASCIIOnly)
	require.Equal(t, "/js", opt.App.StaticFile.SubFolder["js"])
	require.Equal(t, "/static/*", opt.App.RoutePattern.StaticPage)

	t.Log(opt.App.Template)
}

var envfile string = `
COMPOSE_PROJECT_NAME=news-sentiment-analyzer

APP_PORT=8000
APP_HOST=127.0.0.1
APP_API_VERSION=v1

POSTGRES_USERNAME=admin
POSTGRES_HOST=localhost
POSTGRES_PORT=5434
POSTGRES_DB_NAME=nsa
POSTGRES_SSL_MODE=disable
`

func TestReadEnvVar(t *testing.T) {
	envVars := strings.Split(envfile, "\n")
	envVarName := []string{}
	for _, ev := range envVars {
		ev = strings.Trim(ev, "")
		if len(ev) > 0 {
			nvPair := strings.Split(ev, "=")
			_ = os.Setenv(nvPair[0], nvPair[1])
			envVarName = append(envVarName, nvPair[0])
		}
	}

	err := viper.BindEnv(envVarName...)
	require.NoError(t, err)
	// env
	t.Log(t, "v1", viper.GetString("APP_API_VERSION"))
	t.Log(t, "admin", viper.GetString("POSTGRES_USERNAME"))
	t.Log(t, "127.0.0.1", viper.GetString("APP_HOST"))
	t.Log(t, 8000, viper.GetInt("APP_HOST"))

	// default value
	t.Log(t, "app", viper.GetString("APP_NAME"))
	t.Log(t, "dev", viper.GetInt("APP_STATE"))
}

// func TestReadOption(t *testing.T) {
// 	err := global.ReadAppVar(
// 		"../secret.json",
// 		"../config/option.json",
// 		"../config/endpoint.json",
// 	)
// 	require.NoError(t, err)

// 	opt := global.AppVar
// 	err = opt.TokenMaker.UpdateSecret()
// 	require.NoError(t, err)
// 	opt.TokenMaker.SignMethod = jwt.SigningMethodHS512

// 	oldSrc := opt.TokenMaker.GetHexSecretString()
// 	err = opt.TokenMaker.UpdateSecret()
// 	require.NoError(t, err)
// 	newSrc := opt.TokenMaker.GetHexSecretString()
// 	require.NotEqual(t, oldSrc, newSrc)

// 	t.Log(global.AppVar)
// }

package global_test

import (
	"bytes"
	"encoding/json"
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

func TestReadMicroservice(t *testing.T) {
	text := `{
      "name": "language-detector",
      "type": "gRPC",
      "host": "localhost",
      "port": 50051
    }`

	viper.SetConfigType("json")
	err := viper.ReadConfig(strings.NewReader(text))
	require.NoError(t, err)

	var ms1 global.Microservice
	err = viper.Unmarshal(&ms1)
	require.NoError(t, err)

	require.Equal(t, "language-detector", ms1.Name)
	require.Equal(t, "localhost", ms1.Host)
	require.Equal(t, 50051, ms1.Port)
	require.Equal(t, global.APITypegRPC, ms1.Type)
	require.Nil(t, ms1.Params)

	b, err := json.Marshal(ms1)
	require.NoError(t, err)
	require.NotNil(t, b)

	var ms2 global.Microservice
	err = json.Unmarshal(b, &ms2)
	require.NoError(t, err)

	require.Equal(t, ms1.Host, ms2.Host)
	require.Equal(t, ms1.Port, ms2.Port)
	require.Equal(t, ms1.Name, ms2.Name)
	require.Equal(t, ms1.Type, ms2.Type)
}

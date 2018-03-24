package secrets

import (
	"encoding/json"
	"io/ioutil"

	"github.com/romana/rlog"
)

type Secrets struct {
	AppId       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	RedirectUrl string `json:"redirect_url"`
	SentryToken string `json:"sentry_auth_token"`
	SentryDSN   string `json:"sentry_dsn"`
}

type SecretsManager struct {
	secrets Secrets
}

func getSecrets(path string) Secrets {
	var secrets Secrets

	file, err := ioutil.ReadFile(path)
	if err != nil {
		rlog.Error(err)
		return Secrets{}
	}

	if err := json.Unmarshal(file, &secrets); err != nil {
		rlog.Error(err)
		return Secrets{}
	}
	rlog.Debug("Successfully loaded secrets")
	return secrets
}

var secrets_manager *SecretsManager

func LoadSecrets(path string) Secrets {
	secrets_manager = &SecretsManager{getSecrets(path)}
	return secrets_manager.secrets
}

// TODO: concurrency bug
func GetSecrets() Secrets {
	if secrets_manager == nil {
		return Secrets{}
	}
	return secrets_manager.secrets
}

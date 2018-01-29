package secrets

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Secrets struct {
	AppId       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	RedirectUrl string `json:"redirect_url"`
}

type SecretsManager struct {
	secrets Secrets
}

func getSecrets(path string) Secrets {
	var secrets Secrets

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(file, &secrets); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully loaded secrets")
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

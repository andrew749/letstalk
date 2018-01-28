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

func getSecrets() Secrets {
	var secrets Secrets

	file, err := ioutil.ReadFile("secrets/secrets.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(file, &secrets); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully loaded secrets")
	return secrets
}

// TODO: concurrency bug
var secrets_manager *SecretsManager

func GetSecrets() Secrets {
	if secrets_manager == nil {
		secrets_manager = &SecretsManager{getSecrets()}
	}
	return secrets_manager.secrets
}

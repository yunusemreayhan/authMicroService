package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const defaultTokenTimeoutSeconds = 36000

type AppConfig struct {
	TokenTimeoutSeconds int32 `json:"token_timeout_seconds" default:"36000" env:"TOKEN_TIMEOUT_SECONDS"`
}

const configFileName = "/config.json"

func tryReadingConfigFromJSON() (AppConfig, error) {
	// try parsing request as json
	// check if file exist
	err, _ := os.Stat(configFileName)
	if err != nil {
		return AppConfig{}, fmt.Errorf("file %s does not exist", configFileName)
	}

	file, erropen := os.Open(configFileName)
	if erropen != nil {
		return AppConfig{}, fmt.Errorf("can not open file %s", configFileName)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Default().Printf("file.Close error : [%v]\n", err)
		}
	}(file)

	var appConfigRet AppConfig
	var bytesOfJsonFile []byte
	_, errFileRead := file.Read(bytesOfJsonFile)

	if errFileRead != nil {
		return AppConfig{}, fmt.Errorf("file %s does not exist", configFileName)
	}
	errUnmarshal := json.Unmarshal(bytesOfJsonFile, &appConfigRet)

	if errUnmarshal != nil {
		return AppConfig{}, fmt.Errorf("was not able to unmarshal file content [%s] into struct", configFileName)
	}
	return appConfigRet, nil
}

func tryReadingConfigFromEnv() (AppConfig, error) {
	// try parsing request as json
	var appConfigRet AppConfig
	return appConfigRet, nil
}

func GetConfig() AppConfig {
	// try parsing request as json
	var config AppConfig
	config, err := tryReadingConfigFromJSON()
	if err != nil {
		log.Printf("failed to parse json : [%v]", err)
		config, err = tryReadingConfigFromEnv()
		if err != nil {
			log.Printf("failed to parse env : [%v]", err)
			config = AppConfig{
				TokenTimeoutSeconds: defaultTokenTimeoutSeconds,
			}
		}
	}
	if config.TokenTimeoutSeconds == 0 {
		config.TokenTimeoutSeconds = defaultTokenTimeoutSeconds
	}
	log.Default().Printf("config : [%v]", config)
	return config
}

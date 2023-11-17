package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const defaultTokenTimeoutSeconds = 36000

type AppConfig struct {
	TokenTimeoutSeconds int32  `json:"token_timeout_seconds" default:"36000" env:"TOKEN_TIMEOUT_SECONDS"`
	DefaultKeyPath      string `json:"default_key_path" default:"/private.key" env:"DEFAULT_KEY_PATH"`
}

const configFileName = "/config.json"

var config *AppConfig

func tryReadingConfigFromJSON() (AppConfig, error) {
	if config != nil {
		return *config, nil
	}

	config = new(AppConfig)

	// try parsing request as json
	// check if file exist
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

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return AppConfig{}, fmt.Errorf("stat call error, file %s does not exist", configFileName)
	}
	fileSize := fileInfo.Size()

	// Read the entire file into a byte slice
	bytesOfJsonFile := make([]byte, fileSize)
	_, errFileRead := file.Read(bytesOfJsonFile)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file:", err)
		return AppConfig{}, fmt.Errorf("read into byte array error, file %s does not exist", configFileName)
	}

	// Use the 'data' byte slice as needed
	fmt.Printf("File content as byte slice:\n%s\n", string(bytesOfJsonFile))

	if errFileRead != nil {
		return AppConfig{}, fmt.Errorf("file %s does not exist", configFileName)
	}

	errUnmarshal := json.Unmarshal(bytesOfJsonFile, config)

	if errUnmarshal != nil {
		return AppConfig{}, fmt.Errorf("was not able to unmarshal file content [%s] into struct [%v] error [%v]", string(bytesOfJsonFile), *config, errUnmarshal)
	}
	return *config, nil
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
	if config.DefaultKeyPath == "" {
		config.DefaultKeyPath = "/private.key"
	}
	log.Default().Printf("config : [%v]", config)
	return config
}

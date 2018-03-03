package gi

import (
	"encoding/json"
	"os"
	"io/ioutil"
	"fmt"
)

type Config struct {
	Version     int    `json:"version"`
	GithubToken string `json:"github_token"`
}

const configPath = "./config.json"

func (config *Config) init() {
	config.Version = 1
	config.GithubToken = ""
}

func (config *Config) Save() {
	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("Could not save configs")
	}
	data, _ := json.Marshal(&config)
	f.WriteString(string(data))
}

func GetConfig() *Config {
	config := Config{}
	config.Get()
	return &config
}

func (config *Config) Get() bool {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		config.init()
		return false
	}
	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		config.init()
		return false
	}
	json.Unmarshal(raw, &config)
	if config.Version == 0 {
		config.init()
		return false
	}
	return true
}

package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Web3     Web3Config     `json:"web3"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Env  string `json:"env"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	MaxIdle  int    `json:"max_idle"`
	MaxOpen  int    `json:"max_open"`
}

type Web3Config struct {
	RPCURL      string   `json:"rpc_url"`
	ContractAddr string `json:"contract_addr"`
	StartBlock  uint64   `json:"start_block"`
}

var GlobalConfig *Config

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	GlobalConfig = config
	return config, nil
}

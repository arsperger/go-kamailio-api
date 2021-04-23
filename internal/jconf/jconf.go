package jconf

import (
	"encoding/json"
	"os"
)

// ServiceConfig is struct for the config
type ServiceConfig struct {
	Kamailio struct {
		ServerAddr string `json:"jsonrpcs_address"`
	} `json:"kamailio"`
	HTTP struct {
		ListenAddr string `json:"listen_address"`
	} `json:"http"`
	DB struct {
		DbURL string `json:"dburl"`
	} `json:"database"`
}

// LoadConfigFile loads and parse confg
func (s *ServiceConfig) LoadConfigFile(configFile string) error {

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&s); err != nil {
		return err
	}

	return nil
}

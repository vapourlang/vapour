package config

import (
	"encoding/json"
	"os"
	"path"
)

var file string = ".vapour"

type lspConfig struct {
	When     []string `json:"when"`
	Severity []string `json:"severity"`
}

type libPath string

type Config struct {
	Lsp     *lspConfig `json:"lsp"`
	Library libPath    `json:"library"`
}

func makeConfigPath(conf string) string {
	dirname, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return path.Join(dirname, conf)
}

func hasConfig(conf string) bool {
	p := makeConfigPath(conf)

	if p == "" {
		return false
	}

	_, err := os.Stat(p)

	return !os.IsNotExist(err)
}

func ReadConfig() *Config {
	configuration := &Config{
		Library: "",
		Lsp: &lspConfig{
			When:     []string{"open", "save", "close", "text"},
			Severity: []string{"fatal", "warn", "info", "hint"},
		},
	}

	if !hasConfig(file) {
		return configuration
	}

	p := makeConfigPath(file)
	data, err := os.ReadFile(p)

	if err != nil {
		return configuration
	}

	err = json.Unmarshal(data, configuration)

	if err != nil {
		return configuration
	}

	return configuration
}

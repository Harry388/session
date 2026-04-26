package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Location     string
	SearchPaths  []string `json:"searchPaths"`
	IncludePaths []string `json:"includePaths"`
}

func ParseFromConfigDir(location string) (Config, error) {
	file, err := os.Open(filepath.Join(location, "config.json"))
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	conf, err := parseConfig(file)
	conf.Location = location
	return conf, err
}

func parseConfig(r io.Reader) (Config, error) {
	var conf Config
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&conf)
	return conf, err
}

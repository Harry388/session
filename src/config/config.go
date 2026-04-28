package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Location     string
	SearchPaths  []string `json:"searchPaths"`
	IncludePaths []string `json:"includePaths"`
}

const defaultConfigContent = `{
	"searchPaths": [
	],
	"includePaths": [
	]
}`

func ParseFromConfigDir(location string) (Config, error) {
	file, err := os.Open(filepath.Join(location, "config.json"))
	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(location, 0755) // u: rwx, g: r-x, o: r-x
		if err != nil {
			return Config{}, err
		}
		file, err = os.OpenFile(filepath.Join(location, "config.json"), os.O_CREATE|os.O_WRONLY, 0644) // u: rw-, g: r--, o: r--
		if err != nil {
			return Config{}, err
		}
		defer file.Close()
		_, err = file.WriteString(defaultConfigContent)
		if err != nil {
			return Config{}, err
		}
		return parseConfig(strings.NewReader(defaultConfigContent))
	}
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

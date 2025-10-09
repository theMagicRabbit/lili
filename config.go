package main

import (
	"os"
	"path"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Directories map[string]string
}

func (s *State) ReadConfig() error {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	
	liliConfigFile := path.Join(userConfDir, s.ConfigFileName)
	liliConfigData, err := os.ReadFile(liliConfigFile)
	if err != nil {
		return err
	}
	var c Config
	err = toml.Unmarshal(liliConfigData, &c)
	if err != nil {
		return err
	}
	s.Config = &c
	return nil
}


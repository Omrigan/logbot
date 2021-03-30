package main

import (
	"github.com/omrigan/logbot/pkg/storage"
	"os"

	"gopkg.in/yaml.v2"
)

type RecordType struct {
	Options []string `yaml:"options"`
	Comment bool     `yaml:"comment"`
}

type Config struct {
	RecordTypes map[string]*RecordType `yaml:"record_types"`
	Dir         string                 `yaml:"dir"`
	Token       string                 `yaml:"token"`
	Influx      *storage.InfluxConfig  `yaml:"influx"`
	File        *storage.FileConfig    `yaml:"file"`
}

func readConfig() (*Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

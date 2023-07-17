package commons

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigEntry struct {
	Username  string `yaml:"USERNAME"`
	Password  string `yaml:"PASSWORD"`
	Cluster   string `yaml:"CLUSTER"`
	Port      int    `yaml:"PORT"`
	Database  string `yaml:"DATABASE"`
	Schema    string `yaml:"SCHEMA"`
	TableUser string `yaml:"TABLE_USER"`
}

type Config map[string]ConfigEntry

func ReadConfig(file string) (Config, error) {
	var config Config
	data, err := os.ReadFile(file)
	if err != nil {
		log.Printf("cannot read config \nreason:\n%s", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("cannot unmarshal config \nreason:\n%s", err)
	}

	return config, err
}

func (c *Config) ConnString(dbType string) string {
	ce := (*c)[dbType]
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d", ce.Username, ce.Password, ce.Cluster, ce.Port)
	return connString
}

func (c *Config) TableName(dbType string) string {
	ce := (*c)[dbType]
	fullTableName := fmt.Sprintf("%s.%s.%s", ce.Database, ce.Schema, ce.TableUser)
	return fullTableName
}

package commons

import (
	"fmt"
	"log"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/linkedin"
	"gopkg.in/yaml.v2"
)

type ConfigEntry struct {
	Username           string `yaml:"USERNAME"`
	Password           string `yaml:"PASSWORD"`
	Cluster            string `yaml:"CLUSTER"`
	Port               int    `yaml:"PORT"`
	Database           string `yaml:"DATABASE"`
	Schema             string `yaml:"SCHEMA"`
	TableUser          string `yaml:"TABLE_USER"`
	TablePasswordStore string `yaml:"TABLE_PASSWORDSTORE"`
}

func (ce *ConfigEntry) GetConnString() string {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d", ce.Username, ce.Password, ce.Cluster, ce.Port)
	return connString
}

func (ce *ConfigEntry) GetUserTableName() string {
	fullTableName := fmt.Sprintf("%s.%s.%s", ce.Database, ce.Schema, ce.TableUser)
	return fullTableName
}

func (ce *ConfigEntry) GetPasswordStoreTableName() string {
	fullTableName := fmt.Sprintf("%s.%s.%s", ce.Database, ce.Schema, ce.TablePasswordStore)
	return fullTableName
}

type Config map[string]*ConfigEntry

func NewConfig(file string) (Config, error) {
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

type SSOEntry struct {
	ClientID string `yaml:"CLIENTID"`
	Secret   string `yaml:"SECRET"`
	Callback string `yaml:"CALLBACK"`
}
type SSOConfig map[string]*SSOEntry

func NewSSOConfig(file string) (SSOConfig, error) {
	var config SSOConfig
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

func (s *SSOConfig) InitProviders() {
	goth.UseProviders(
		google.New((*s)["google"].ClientID, (*s)["google"].Secret, (*s)["google"].Callback, "email", "profile"),
		github.New((*s)["github"].ClientID, (*s)["github"].Secret, (*s)["github"].Callback, "email"),
		linkedin.New((*s)["linkedin"].ClientID, (*s)["linkedin"].Secret, (*s)["linkedin"].Callback),
	)
}

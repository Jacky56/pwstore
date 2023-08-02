package commons

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	got, err := NewConfig("../configs/test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	want := Config{
		"COCKROACH_DB": {
			Username:           "dummy",
			Password:           "dummy",
			Cluster:            "dummy",
			Port:               12345,
			Database:           "pwstore",
			Schema:             "test",
			TableUser:          "user",
			TablePasswordStore: "passwordStore",
		},
	}

	if err != nil {
		t.Errorf("%s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("maps are not the same!\ngot:\n%+v\nwant:\n%+v", got, want)
	}
	fmt.Println(want["COCKROACH_DB"])
	fmt.Println(got["COCKROACH_DB"])
}

func TestConfigConnString(t *testing.T) {
	config, err := NewConfig("../configs/test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	want := "postgresql://dummy:dummy@dummy:12345"
	got := config["COCKROACH_DB"].GetConnString()
	if got != want {
		t.Errorf("connection string not same!\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestConfigTableName(t *testing.T) {
	config, err := NewConfig("../configs/test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	want := "pwstore.test.user"
	got := config["COCKROACH_DB"].GetUserTableName()
	if got != want {
		t.Errorf("connection string not same!\ngot:\n%s\nwant:\n%s", got, want)
	}
	fmt.Println(got)
}

func TestNewSSOConfig(t *testing.T) {
	got, err := NewSSOConfig("../configs/SSO_test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	want := SSOConfig{
		"google": {
			ClientID: "google",
			Secret:   "google",
			Callback: "http://localhost:4000/auth/google/callback",
		},
		"github": {
			ClientID: "github",
			Secret:   "github",
			Callback: "http://localhost:4000/auth/github/callback",
		},
	}

	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(want["google"], want["github"])
	fmt.Println(got["google"], got["github"])
	if !reflect.DeepEqual(got, want) {
		t.Errorf("maps are not the same!\ngot:\n%+v\nwant:\n%+v", got, want)
	}
}

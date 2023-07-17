package commons

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	got, err := ReadConfig("../configs/test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	want := Config{
		"COCKROACH_DB": {
			Username:  "dummy",
			Password:  "dummy",
			Database:  "pwstore",
			Schema:    "test",
			TableUser: "user",
			Port:      12345,
			Cluster:   "dummy",
		},
	}

	if err != nil {
		t.Errorf("%s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("maps are not the same!\ngot:\n%+v\nwant:\n%+v", got, want)
	}
	fmt.Printf("%+v", want["COCKROACH_DB"])
}

func TestConfigConnString(t *testing.T) {
	config, err := ReadConfig("../configs/test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	want := "postgresql://dummy:dummy@dummy:12345"
	got := config.ConnString("COCKROACH_DB")
	if got != want {
		t.Errorf("connection string not same!\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestConfigTableName(t *testing.T) {
	config, err := ReadConfig("../configs/test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	want := "pwstore.test.user"
	got := config.TableName("COCKROACH_DB")
	if got != want {
		t.Errorf("connection string not same!\ngot:\n%s\nwant:\n%s", got, want)
	}
}

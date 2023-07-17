package databases

import (
	"pwstore/commons"
	"testing"
)

func TestCockroachdbConnectFail(t *testing.T) {
	config, err := commons.ReadConfig("../configs/test.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	db := NewCockroachdb[any](config.ConnString("COCKROACH_DB"), config.TableName("COCKROACH_DB"))
	_, err = db.Connect()
	if err == nil {
		t.Errorf("helloworld is not a real connection string")
	}

	conn := db.GetConn()
	if conn != nil {
		t.Errorf("%+v", conn)
	}
}

func TestCockroachdbConnectPass(t *testing.T) {
	config, err := commons.ReadConfig("../configs/dev.yaml")
	if err != nil {
		t.Errorf("cannot load config!\n%s", err)
	}
	db := NewCockroachdb[any](config.ConnString("COCKROACH_DB"), config.TableName("COCKROACH_DB"))
	_, err = db.Connect()
	if err != nil {
		t.Errorf("helloworld is not a real connection string")
	}
	defer db.Close()

	conn := db.GetConn()
	if conn == nil {
		t.Errorf("%+v", conn)
	}
}

func TestCockroachdbQueryPass(t *testing.T) {

}
func TestCockroachdbQueryFail(t *testing.T) {

}

package databases

import (
	"fmt"
	"pwstore/commons"
	"pwstore/data"
	"pwstore/types"
	"testing"
)

func setup(conf string) types.Database {
	config, err := commons.NewConfig(conf)
	if err != nil {
		fmt.Printf("cannot load config!\n%s", err)
	}
	db := NewCockroachdb(config["COCKROACH_DB"])
	_, err = db.Connect()
	return db
}

func TestCockroachdbConnectFail(t *testing.T) {
	db := setup("../configs/test.yaml")

	conn := db.GetConn()
	if conn != nil {
		t.Errorf("%+v", conn)
	}
}

func TestCockroachdbConnectPass(t *testing.T) {
	db := setup("../configs/dev.yaml")
	defer db.Close()

	conn := db.GetConn()
	if conn == nil {
		t.Errorf("%+v", conn)
	}
}

func TestCockroachdbQueryPass(t *testing.T) {
	db := setup("../configs/dev.yaml")
	config, err := commons.NewConfig("../configs/dev.yaml")
	defer db.Close()

	user, err := db.Query(
		fmt.Sprintf("select * from %s", config["COCKROACH_DB"].GetUserTableName()),
	)
	if err != nil {
		t.Errorf("cannot query user table %s", err)
	}
	fmt.Printf("%+v", user)
}

func TestCockroachdbUserPass(t *testing.T) {
	db := setup("../configs/dev.yaml")
	defer db.Close()

	user, err := db.GetUser("jackybanh1997@gmail.com")
	if err != nil {
		t.Errorf("cannot query user table %s", err)
	}
	fmt.Printf("%+v", user)
}

func TestCockroachdbPasswordsPass(t *testing.T) {
	db := setup("../configs/dev.yaml")
	defer db.Close()

	user, err := db.GetUser("jacky@gmail.com")
	if err != nil {
		t.Errorf("cannot query user table %s", err)
	}

	passwordStore, err := db.GetPasswordStore(user.Uuid)
	if err != nil {
		t.Errorf("cannot query passwordstore table %s", err)
	}
	fmt.Printf("%+v", passwordStore)
}

func TestSetUser(t *testing.T) {
	db := setup("../configs/dev.yaml")
	defer db.Close()
	user := data.User{
		Email:    "foo@bar.com",
		Password: "jacky1",
	}
	err := db.SetUser(&user)
	if err != nil {
		t.Errorf("cannot upsert\n%s", err)
	}
}

func TestListUsers(t *testing.T) {
	db := setup("../configs/dev.yaml")
	defer db.Close()
	users, err := db.ListUsers()
	if err != nil {
		t.Errorf("cannot get list of users\n%s", err)
	}
	fmt.Println(users)
}

func TestSetPasswordStore(t *testing.T) {
	db := setup("../configs/dev.yaml")
	defer db.Close()
	user, _ := db.GetUser("maz@bar.com")

	pws := data.PasswordStore{
		Uuid: user.Uuid,
		PasswordStore: map[string]string{
			"hello": "world",
			"log":   "fer",
			"yo":    "man",
		},
	}

	err := db.SetPasswordStore(&pws)
	if err != nil {
		t.Errorf("cannot set password!\n%s", err)
	}
}

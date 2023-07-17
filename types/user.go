package types

type User struct {
	Uuid          string            `json:"uuid"`
	Username      string            `json:"username"`
	Password      string            `json:"password"`
	PasswordStore map[string]string `json:"Password_store"`
}

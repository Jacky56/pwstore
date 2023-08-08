package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"net/http"
	"net/http/httptest"
	"pwstore/commons"
	"pwstore/databases"
	"pwstore/types"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func setupDB() *types.Database {
	config, _ := commons.NewConfig("../configs/dev.yaml")
	db := databases.NewCockroachdb(config["COCKROACH_DB"])
	return &db
}

func setupPostBody(body any) *bytes.Reader {
	postBody, _ := json.Marshal(body)
	responseBody := bytes.NewReader(postBody)
	return responseBody
}

func decodeResponse(resp *http.Response) map[string]any {
	var m map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()
	return m
}

func TestCreateJWT(t *testing.T) {
	id, _ := uuid.NewUUID()
	claims := jwt.MapClaims{
		"uuid":  id,
		"email": "hello@world.com",
		"exp":   1233243243254254325, // time.Now().Add(time.Second * 3).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	singedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token)
	fmt.Println(singedToken, len(singedToken))
}

func TestDecryptJWT(t *testing.T) {
	// jwtConfig := jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{Key: []byte("dsadsa")},
	// }

	// keyFunc := func(token *jwt.Token) (interface{}, error) {
	// 	return jwtConfig.SigningKey.Key, nil
	// }

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	}

	//singedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhlbGxvQHdvcmxkLmNvbSIsImV4cCI6MTIzMzI0MzI0MzI1NDI1NDMyNSwidXVpZCI6IjJjOTc2NjcxLTJhNDctMTFlZS05MDFjLTAwMTU1ZDVjOGJlMSJ9.fECHRWO0TmNnjrSogc84O_weTtgqbocZ6cPX0a3lMtA"
	singedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImphY2t5QGdhbWlsLmNvbSIsImV4cCI6MTY5MDI0Nzc1NCwidXVpZCI6ImRmMjY4NzkzLTQ4MjgtNDVhZC1hOTc3LWNjMzY1ZTRmOTkwZSJ9.9qtdbTzLbuC6KipNyjPtlJiL5JracB7QDw_Fml_NQX0"

	token, err := jwt.Parse(singedToken, keyFunc)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token.Valid)
	fmt.Println(token.Claims)

}

func TestAuth(t *testing.T) {
	db := setupDB()
	(*db).Connect()
	app := fiber.New()
	auth := NewAuthentication(db)
	app.Post("/auth", auth.Auth)

	postBody := setupPostBody(map[string]any{
		"email":    "jacky@gamil.com",
		"password": "bannh",
	})

	newreq := httptest.NewRequest("POST", "/auth", postBody)
	newreq.Header.Set("Content-Type", "application/json")

	req, err := app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}

	m := decodeResponse(req)
	fmt.Println(m["token"])

}

func TestCreateUserJWT(t *testing.T) {
	db := setupDB()
	(*db).Connect()
	app := fiber.New()
	auth := NewAuthentication(db)
	app.Get("/create/:token", auth.CreateUserJWT)

	claims := jwt.MapClaims{
		"email":    "hello@world.com",
		"password": "helloworld",
		"exp":      1233243243254254325, // time.Now().Add(time.Second * 3).Unix(),
	}

	token := jwt.NewWithClaims(auth.SigningMethod, claims)
	singedToken, _ := token.SignedString(auth.signedString)

	tokenGet := fmt.Sprintf("/create/%s", singedToken)

	newreq := httptest.NewRequest("GET", tokenGet, nil)

	req, err := app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}

	// m := decodeResponse(req)
	fmt.Println(req.Body)

}
func TestCreateUser(t *testing.T) {
	db := setupDB()
	(*db).Connect()
	app := fiber.New()
	auth := NewAuthentication(db)
	app.Post("/create/account", auth.CreateUser)

	postBody := setupPostBody(map[string]any{
		"email":           "maz@bar.com",
		"emailConfirm":    "maz@bar.com",
		"password":        "mazbar",
		"passwordConfirm": "mazbar",
	})

	newreq := httptest.NewRequest("POST", "/create/account", postBody)
	newreq.Header.Set("Content-Type", "application/json")

	req, err := app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}

	// m := decodeResponse(req)
	t.Log(req.Body)

	t.Log((*db).ListUsers())
}

func TestEmailCreateUserLink(t *testing.T) {
	db := setupDB()
	(*db).Connect()
	app := fiber.New()
	auth := NewAuthentication(db)
	app.Post("/register/create", auth.EmailCreateUserJWTLink)

	postBody := setupPostBody(map[string]any{
		"email":    "@bar.com",
		"password": "faz",
	})

	newreq := httptest.NewRequest("POST", "/register/create", postBody)
	newreq.Header.Set("Content-Type", "application/json")

	req, err := app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Println(req.Body)

}

func TestSHA256(t *testing.T) {
	h := sha256.New()
	s := []byte(fmt.Sprintf("%d", time.Now().Unix()))
	h.Write(s)
	bs := h.Sum(nil)
	fmt.Println(s)
	fmt.Printf("%x\n", bs)
}

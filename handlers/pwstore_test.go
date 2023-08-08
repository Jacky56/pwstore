package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSetPW(t *testing.T) {
	db := setupDB()
	(*db).Connect()
	app := fiber.New()
	auth := NewAuthentication(db)
	pwstore := NewPwStore(db, auth)
	app.Post("/auth", auth.Auth)
	app.Post("/updatepws", auth.GetSignedToken, pwstore.UpdatePw)

	u, err := (*db).GetUser("maz@bar.com")
	if err != nil {
		t.Error(err)
	}
	d := map[string]any{
		"email":    u.Email,
		"password": u.Password,
	}
	postBody := setupPostBody(d)
	newreq := httptest.NewRequest("POST", "/auth", postBody)
	newreq.Header.Set("Content-Type", "application/json")
	req, err := app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}

	log.Println(req.Cookies()[0])

	postString := bytes.NewReader([]byte("keyword=hello&value=reflect&keyword=use&value=it"))
	newreq = httptest.NewRequest("POST", "/updatepws", postString)
	newreq.AddCookie(
		&http.Cookie{
			Name:  req.Cookies()[0].Name,
			Value: req.Cookies()[0].Value,
		},
	)
	newreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req, err = app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(req.Body)
}

func TestGetPW(t *testing.T) {
	db := setupDB()
	(*db).Connect()
	app := fiber.New()
	auth := NewAuthentication(db)
	pwstore := NewPwStore(db, auth)
	app.Post("/auth", auth.Auth)
	app.Get("/pw", auth.AuthAllow, pwstore.GetPW)

	u, err := (*db).GetUser("maz@bar.com")
	if err != nil {
		t.Error(err)
	}
	d := map[string]any{
		"email":    u.Email,
		"password": u.Password,
	}
	postBody := setupPostBody(d)
	newreq := httptest.NewRequest("POST", "/auth", postBody)
	newreq.Header.Set("Content-Type", "application/json")
	req, err := app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}

	newreq = httptest.NewRequest("GET", "/pw", postBody)

	newreq.Header.Set("Content-Type", "application/json")
	newreq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.Cookies()[0].Value))
	req, err = app.Test(newreq, -1)
	if err != nil {
		t.Errorf("%s", err)
	}
	s, _ := io.ReadAll(req.Body)
	var pwsD map[string]any

	want := map[string]any{
		"hello": "reflect",
		"use":   "it",
	}
	json.Unmarshal(s, &pwsD)
	check := reflect.DeepEqual(pwsD, want)
	if !check {
		t.Error("not the same dicts", pwsD, want)
	}
}

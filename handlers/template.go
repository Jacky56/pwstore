package handlers

import (
	"fmt"
	"pwstore/data"
	"pwstore/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/django/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Template struct {
	Engine *django.Engine
	db     *types.Database
	auth   *Authentication
}

func NewTemplate(db *types.Database, auth *Authentication) *Template {
	return &Template{
		db:     db,
		Engine: django.New("./templates", ".html"),
		auth:   auth,
	}
}

func (t *Template) Partial(c *fiber.Ctx) error {
	partial := c.Params("partial")
	err := c.Render(fmt.Sprintf("partials/%s", partial), nil)
	if err != nil {
		log.Warn(fmt.Sprintf("%s partial broken", partial), err)
		return err
	}
	return nil
}

func (t *Template) Login(c *fiber.Ctx) error {
	return c.Render("partials/login", nil, "index")
}

func (t *Template) Register(c *fiber.Ctx) error {
	return c.Render("partials/register", nil, "index")
}

func (t *Template) Wait(c *fiber.Ctx) error {
	time.Sleep(time.Second * 2)
	return c.Next()
}

func (t *Template) Index(c *fiber.Ctx) error {
	signedToken := c.Locals(t.auth.LocalsName)
	c.Set("HX-Refresh", "true")
	if signedToken == nil {
		log.Info("User not logged in redirect to login")
		return c.Redirect("/login")
	}

	claims := signedToken.(*jwt.Token).Claims.(jwt.MapClaims)

	id, err := uuid.Parse(claims["uuid"].(string))
	if err != nil {
		log.Warn("UUID unparsable, redirect to login")
		return c.Redirect("/login")
	}
	u := data.User{
		Email:    claims["email"].(string),
		Provider: claims["provider"].(string),
		Uuid:     id,
	}
	pwstore, err := (*t.db).GetPasswordStore(u.Uuid)
	if err != nil {
		log.Warn(err)
	}
	return c.Render("partials/pwstore", fiber.Map{
		"user": u,
		"pws":  pwstore.PasswordStore,
	}, "index")

}

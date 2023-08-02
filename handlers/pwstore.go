package handlers

import (
	"net/url"
	"pwstore/commons"
	"pwstore/data"
	"pwstore/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PwStore struct {
	db   *types.Database
	auth *Authentication
}

func NewPwStore(db *types.Database, auth *Authentication) *PwStore {
	return &PwStore{
		db:   db,
		auth: auth,
	}
}

func (p *PwStore) UpdatePw(c *fiber.Ctx) error {

	signedToken := c.Locals(p.auth.LocalsName)
	if signedToken == nil {
		log.Info("User not logged in redirect to login")
		return c.Redirect("/login")
	}
	claims := signedToken.(*jwt.Token).Claims.(jwt.MapClaims)
	uuid, err := uuid.Parse(claims["uuid"].(string))
	if err != nil {
		log.Warn("UUID unparsable, redirect to login")
		return c.Redirect("/login")
	}

	v, err := url.ParseQuery(string(c.Body()))
	if err != nil {
		s := "cannot parse pws query"
		log.Warn(s)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert(s),
		})
	}

	d := make(map[string]string)
	for i, k := range v["keyword"] {
		if len(k) == 0 {
			continue
		}
		d[k] = v["value"][i]
	}
	pws := &data.PasswordStore{
		Uuid:          uuid,
		PasswordStore: d,
	}

	err = (*p.db).SetPasswordStore(pws)
	if err != nil {
		s := "cannot set password stores"
		log.Warn(s)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert(s),
		})
	}
	c.SendStatus(fiber.StatusAccepted)
	c.Set("HX-Refresh", "true")

	return c.Next()
}

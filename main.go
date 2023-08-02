package main

import (
	"log"
	"pwstore/commons"
	"pwstore/data"
	"pwstore/databases"
	"pwstore/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/shareed2k/goth_fiber"
)

func main() {

	config, _ := commons.NewConfig("./configs/dev.yaml")
	sso, _ := commons.NewSSOConfig("./configs/SSO_dev.yaml")
	db := databases.NewCockroachdb(config["COCKROACH_DB"])
	sso.InitProviders()

	db.Connect()
	defer db.Close()
	auth := handlers.NewAuthentication(&db)
	template := handlers.NewTemplate(&db, auth)
	pwstore := handlers.NewPwStore(&db, auth)
	app := fiber.New(fiber.Config{
		PassLocalsToViews: true,
		Views:             template.Engine,
	})

	app.Static("/statics", "./statics")

	app.Get("/", auth.GetSignedToken, template.Index)

	app.Get("/register", auth.AuthDeny, template.Register)
	app.Post("/register", auth.AuthDeny, auth.CreateUser)
	app.Get("/login", auth.AuthDeny, template.Login)
	app.Get("/logout", auth.Logout)
	app.Get("/partials/:partial", template.Partial)

	app.Get("/login/:provider", auth.AuthDeny, goth_fiber.BeginAuthHandler)
	app.Get("/auth/:provider/callback", auth.AuthSSO)
	app.Post("/auth", auth.AuthDeny, template.Wait, auth.Auth)

	// app.Post("/register", auth.AuthDeny, auth.EmailCreateUserJWTLink)

	// app.Get("/register/:token", auth.AuthDeny, auth.CreateUserJWT)

	// app.Get("/login", auth.AuthDeny, auth.Login)

	app.Get("/restricted", auth.AuthAllow, template.Restricted)
	app.Post("/updatepws", auth.AuthAllow, auth.GetSignedToken, template.Wait, pwstore.UpdatePw)

	log.Fatal(app.Listen(":4000"))
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	u := data.User{
		Email: claims["email"].(string),
	}
	return c.SendString("Welcome " + u.Email)
}

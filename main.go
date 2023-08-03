package main

import (
	"flag"
	"fmt"
	"log"
	"pwstore/commons"
	"pwstore/databases"
	"pwstore/handlers"

	"github.com/gofiber/fiber/v2"

	"github.com/shareed2k/goth_fiber"
)

func main() {
	env := flag.String("env", "dev", "Which configs to use, default: dev")
	port := flag.String("port", "8080", "Which configs to use, default: 8080")
	flag.Parse()
	fmt.Printf("env: %s, port: %s", *env, *port)
	config, _ := commons.NewConfig(fmt.Sprintf("./configs/%s.yaml", *env))
	sso, _ := commons.NewSSOConfig(fmt.Sprintf("./configs/SSO_%s.yaml", *env))
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
	app.Post("/register", auth.AuthDeny, template.Wait, auth.CreateUser, auth.Auth)
	app.Get("/login", auth.AuthDeny, template.Login)
	app.Get("/logout", auth.Logout)
	app.Get("/partials/:partial", template.Partial)

	app.Get("/login/:provider", auth.AuthDeny, goth_fiber.BeginAuthHandler)
	app.Get("/auth/:provider/callback", auth.AuthSSO)
	app.Post("/auth", auth.AuthDeny, template.Wait, auth.Auth)

	app.Get("/restricted", auth.AuthAllow, template.Restricted)
	app.Post("/updatepws", auth.AuthAllow, auth.GetSignedToken, template.Wait, pwstore.UpdatePw)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", *port)))
}

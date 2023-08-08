package handlers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"pwstore/commons"
	"pwstore/data"
	"pwstore/types"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shareed2k/goth_fiber"
)

func Hellworld() {
	fmt.Println("yes")
}

type Authentication struct {
	db            *types.Database
	signedString  []byte
	expiry        int64
	SigningMethod jwt.SigningMethod
	CookieName    string
	LocalsName    string
	AuthAllow     fiber.Handler
	AuthDeny      fiber.Handler
	keyFunc       func(*jwt.Token) (any, error)
}

func NewAuthentication(db *types.Database) *Authentication {
	cookieName := "pwstore-jwt"
	localsName := "local_jwt" // c.Locals are not allowed to have -
	signedString := []byte(fmt.Sprintf("%d", time.Now().Unix()))
	a := Authentication{
		db:            db,
		signedString:  signedString,
		expiry:        time.Now().Add(time.Minute * 30).Unix(),
		SigningMethod: jwt.SigningMethodHS256,
		CookieName:    cookieName,
		LocalsName:    localsName,
		keyFunc: func(token *jwt.Token) (any, error) {
			return signedString, nil
		},
	}

	a.AuthAllow = jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: a.signedString},
		TokenLookup:  fmt.Sprintf("header:Authorization,cookie:%s", cookieName),
		AuthScheme:   "Bearer",
		ErrorHandler: func(c *fiber.Ctx, e error) error { return c.SendFile("templates/login.html") },
	})

	a.AuthDeny = jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: a.signedString},
		TokenLookup:    fmt.Sprintf("header:Authorization,cookie:%s", cookieName),
		AuthScheme:     "Bearer",
		ErrorHandler:   func(c *fiber.Ctx, e error) error { return c.Next() },
		SuccessHandler: func(c *fiber.Ctx) error { return c.Redirect("/") },
	})

	return &a
}

func (a *Authentication) createToken(u *data.User, expiry int64) (string, error) {
	claims := jwt.MapClaims{
		"uuid":     u.Uuid,
		"email":    u.Email,
		"provider": u.Provider,
		"exp":      expiry,
	}
	token := jwt.NewWithClaims(a.SigningMethod, claims)

	t, err := token.SignedString(a.signedString)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (a *Authentication) Auth(c *fiber.Ctx) error {
	var u data.User
	c.BodyParser(&u)

	user, err := (*a.db).GetUser(u.Email)
	if err != nil {
		log.Warn(err, " email: ", u.Email)
		c.SendStatus(fiber.StatusUnauthorized)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert("Incorrect Email or Password"),
		})
	}
	if u.Password != user.Password {
		log.Warn("Password mismatch")
		c.SendStatus(fiber.StatusUnauthorized)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert("Incorrect Email or Password"),
		})
	}

	t, err := a.createToken(&user, a.expiry)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:  a.CookieName,
		Value: t,
	})
	c.Redirect("/")
	c.Set("HX-Refresh", "true")
	return c.JSON(fiber.Map{"token": t})
}

func (a *Authentication) AuthSSO(c *fiber.Ctx) error {
	gothUser, err := goth_fiber.CompleteUserAuth(c)

	log.Info(gothUser.Email)
	if err != nil {
		log.Error(err)
	}
	if len(gothUser.Email) == 0 {
		log.Warn("Email does not exist for SSO")
		c.SendStatus(fiber.StatusInternalServerError)
		text := "<p>Email does not exist for SSO</p>"
		if gothUser.Provider == "github" {
			text += "<p>Please make your github email accessible to the public.</p>"
		}
		return c.Render("partials/login", fiber.Map{
			"alert": commons.NewAlert(text),
		}, "index")
	}

	user, _ := (*a.db).GetUser(gothUser.Email)
	if user == (data.User{}) {
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%d", time.Now().Unix())))
		bs := h.Sum(nil)
		user = data.User{
			Email:    gothUser.Email,
			Provider: gothUser.Provider,
			Password: fmt.Sprintf("%x", bs),
		}
		err := (*a.db).SetUser(&user)
		if err != nil {
			log.Warn("cannot create user")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		// to get uuid for token
		user, _ = (*a.db).GetUser(gothUser.Email)
		if err != nil {
			log.Warn("cannot get user")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	t, err := a.createToken(&user, a.expiry)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:  a.CookieName,
		Value: t,
	})
	c.Redirect("/")
	return c.JSON(fiber.Map{"token": t})
}

func (a *Authentication) createUserMethod(u *data.User) (string, error) {
	users, err := (*a.db).ListUsers()
	if err != nil {
		s := "cannot get list of users when creating account!"
		log.Warn(s, err)
		return s, err
	}
	for _, e := range *users {
		if e.Email == strings.ToLower(u.Email) {
			s := "Email already Exist! email"
			log.Warn(s, u.Email)
			// todo: redirect with message that account was created
			return s, errors.New("Email exist")
		}
	}
	err = (*a.db).SetUser(u)
	if err != nil {
		s := "Failed to create account"
		log.Warn(s, u.Email)
		return s, err
	}
	return "", nil
}

type registerForm struct {
	Email           string `json:"email"`
	EmailConfirm    string `json:"emailConfirm"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func (a *Authentication) CreateUser(c *fiber.Ctx) error {
	var rf registerForm
	err := c.BodyParser(&rf)

	if len(rf.Email) == 0 {
		s := "Email cannot be blank!"
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert(s),
		})
	}

	if rf.Email != rf.EmailConfirm {
		s := "Email does not match!"
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert(s),
		})
	}

	if len(rf.Password) < 6 {
		s := "Make a password longer than 6 characters!"
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert(s),
		})
	}

	if rf.Password != rf.PasswordConfirm {
		s := "Password does not match!"
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert(s),
		})
	}

	var u data.User
	c.BodyParser(&u)
	s, err := a.createUserMethod(&u)
	if err != nil {
		log.Warn("Failed to create accounnt", err)
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("partials/alert", fiber.Map{
			"alert": commons.NewAlert(s),
		})
	}

	return c.Next()
}

// triggered on email get link
func (a *Authentication) CreateUserJWT(c *fiber.Ctx) error {

	jwt_token := c.Params("token")
	token, err := jwt.Parse(jwt_token, a.keyFunc)
	// handles expiry/validity
	if err != nil {
		log.Warnf("CreateUser token fault\n%s", err)
		return c.SendString(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	u := data.User{
		Email:    claims["email"].(string),
		Password: claims["password"].(string),
	}

	// do not need to check for email dupes because of unique constraint

	err = (*a.db).SetUser(&u)
	if err != nil {
		log.Warn("Failed to create account! email: ", u.Email)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendString("account created!")
}

// sends email to user with jwt embed that triggers CreateUserJWT
func (a *Authentication) EmailCreateUserJWTLink(c *fiber.Ctx) error {
	var u data.User
	c.BodyParser(&u)
	users, err := (*a.db).ListUsers()
	if err != nil {
		log.Warn("cannot get list of users when creating account!\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for _, e := range *users {
		if e.Email == strings.ToLower(u.Email) {
			log.Warn("Email already Exist! email: ", u.Email)
			// todo: redirect with message that account was created
			return c.SendStatus(fiber.StatusMethodNotAllowed)
		}
	}

	claims := jwt.MapClaims{
		"email":    u.Email,
		"password": u.Password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(a.SigningMethod, claims)
	singedToken, _ := token.SignedString(a.signedString)

	// todo: url
	link := fmt.Sprintf("/register/%s", singedToken)

	log.Infof("email sent\nemail: %s\nlink %s", u.Email, link)
	return c.Redirect("/")
}

func (a *Authentication) Logout(c *fiber.Ctx) error {
	c.ClearCookie(a.CookieName)
	if err := goth_fiber.Logout(c); err != nil {
		log.Info(err)
	}
	c.Set("HX-Refresh", "true")
	return c.Redirect("/")
}

func (a *Authentication) GetSignedToken(c *fiber.Ctx) error {
	signedToken := c.Cookies(a.CookieName)
	if len(signedToken) == 0 {
		log.Info("Cookies not set")
		return c.Next()
	}
	token, err := jwt.Parse(signedToken, a.keyFunc)
	if err != nil {
		log.Warn(err)
		return c.Next()
	}
	if token.Valid != true {
		log.Warn("token expired", token.Claims)
	}
	c.Locals(a.LocalsName, token)

	return c.Next()
	// if err == nil && token.Valid {
	// 	// Store user information from token into context.
	// 	c.Locals(cfg.ContextKey, token)
	// 	return cfg.SuccessHandler(c)
	// }
	// return cfg.ErrorHandler(c, err)
}

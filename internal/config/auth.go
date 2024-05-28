package config

import (
	"fmt"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
)

type AuthSetup struct {
	cnf    Auth
	google Google
	fb     Facebook
}

func NewAuthSetup(cnf *Config) *AuthSetup {
	return &AuthSetup{cnf: cnf.Auth, google: cnf.Google, fb: cnf.Facebook}
}

func (a *AuthSetup) NewAuth() {
	googleClientID := a.google.ClientID
	googleClientSecret := a.google.ClientSecret
	fbClientID := a.fb.ClientID
	fbClientSecret := a.fb.ClientSecret
	googleCallBackUrl := a.cnf.GoogleCallBackUrl
	fbCallBackUrl := a.cnf.FacebookCallBackUrl

	key := a.cnf.Secret_Key
	email := a.google.ScopeEmail
	profile := a.google.ScopeProfile

	maxAge, _ := strconv.Atoi(a.cnf.MaxAge)
	isProd, _ := strconv.ParseBool(a.cnf.IsProd)

	store := sessions.NewCookieStore([]byte(key))
	if store == nil {
		fmt.Println("Error: Failed to create new cookie store")
		return
	}
	store.MaxAge(maxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, googleCallBackUrl, email, profile),
		facebook.New(fbClientID, fbClientSecret, fbCallBackUrl),
	)

	fmt.Println(goth.GetProviders())
}

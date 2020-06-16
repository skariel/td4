package handlers

// SEE: https://codelike.pro/easy-social-login-oauth-in-go-lang/

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danilopolani/gocialite"
	"github.com/dgrijalva/jwt-go"

	"sql/db"
)

const jwtExpiryDelayHrs = 24 * 12 * 30

// WithGocialInContext self explanatory!
func WithGocialInContext(gocial *gocialite.Dispatcher, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), key(contextKeyGocial), gocial)
	return r.WithContext(ctx)
}

// GetGocialFromContext self explanatory!
func GetGocialFromContext(r *http.Request) *gocialite.Dispatcher {
	return r.Context().Value(key(contextKeyGocial)).(*gocialite.Dispatcher)
}

// SocialRedirectHandler handle login with github
func SocialRedirectHandler(w http.ResponseWriter, r *http.Request) {
	appSettings := map[string]string{
		"clientID":     os.Getenv("TD4_github_client_id"),
		"clientSecret": os.Getenv("TD4_github_client_secret"),
		"redirectURL":  "https://api.solvemytest.dev/auth/github/callback",
	}
	gocial := GetGocialFromContext(r)
	authURL, err := gocial.New().
		Driver("github").
		// Scopes([]string{"user:email"}). included by default by gocialite
		Redirect(
			appSettings["clientID"],
			appSettings["clientSecret"],
			appSettings["redirectURL"],
		)

	if err != nil {
		ise(w, err)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

// SocialCallbackHandler handle login with github
func SocialCallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state := query.Get("state")
	code := query.Get("code")
	q := GetQuerierFromContext(r)
	gocial := GetGocialFromContext(r)

	user, _, err := gocial.Handle(state, code) // token not used
	if err != nil {
		ise(w, err)
		return
	}

	user.ID = "github:" + user.ID
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"display_name": user.Username,
		"email":        user.Email,
		"avatar":       user.Avatar,
		"id":           user.ID,
	})
	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("TD4_JWT_SECRET")))

	if err != nil {
		ise(w, err)
		return
	}

	err = q.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:          user.ID,
		DisplayName: user.Username,
		Email:       user.Email,
		Avatar:      user.Avatar})
	if err != nil {
		ise(w, err)
		return
	}

	jwtCookie := http.Cookie{
		Name:     "jwt_auth",
		Secure:   true,
		Path:     "/",
		Domain:   "solvemytest.dev",
		HttpOnly: false,
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * jwtExpiryDelayHrs),
		SameSite: http.SameSiteStrictMode}
	http.SetCookie(w, &jwtCookie)

	http.Redirect(w, r, os.Getenv("TD4_SOCIAL_AUTH_REDIRECT"), http.StatusFound)
}

// GetUserFromAuthorizationHeader from jwt
func GetUserFromAuthorizationHeader(r *http.Request) *db.Td4User {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Bearer" {
		return nil
	}

	token, err := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TD4_JWT_SECRET")), nil
	})
	if err != nil {
		return nil
	}

	return &db.Td4User{
		DisplayName: token.Claims.(jwt.MapClaims)["display_name"].(string),
		Email:       token.Claims.(jwt.MapClaims)["email"].(string),
		Avatar:      token.Claims.(jwt.MapClaims)["avatar"].(string),
		ID:          token.Claims.(jwt.MapClaims)["id"].(string),
	}
}

// GetUserFromContext self explanatory!
func GetUserFromContext(r *http.Request) *db.Td4User {
	return r.Context().Value(key(contextKeyUser)).(*db.Td4User)
}

// WithUserInContext self explanatory!
func WithUserInContext(user *db.Td4User, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), key(contextKeyUser), user)
	return r.WithContext(ctx)
}

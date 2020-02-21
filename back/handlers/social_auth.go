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

	"../../sql/db"
)

var gocial = gocialite.NewDispatcher()

// SocialRedirectHandler handle login with github
func SocialRedirectHandler(w http.ResponseWriter, r *http.Request) {
	appSettings := map[string]string{
		"clientID":     os.Getenv("github_client_id"),
		"clientSecret": os.Getenv("github_client_secret"),
		"redirectURL":  "http://localhost:8081/auth/github/callback",
	}
	authURL, err := gocial.New().
		Driver("github").
		//Scopes([]string{"user:email"}). included by default by gocialite
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
	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("td4_jwt_secret")))
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
		Secure:   false,
		Path:     "/",
		HttpOnly: false,
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24 * 12 * 30),
		SameSite: http.SameSiteDefaultMode}
	http.SetCookie(w, &jwtCookie)
	userDisplayNameCookie := http.Cookie{
		Name:     "user_display_name",
		Secure:   false,
		Path:     "/",
		HttpOnly: false,
		Value:    user.Username,
		Expires:  time.Now().Add(time.Hour * 24 * 12 * 30),
		SameSite: http.SameSiteDefaultMode}
	http.SetCookie(w, &userDisplayNameCookie)
	userAvatarCookie := http.Cookie{
		Name:     "user_avatar",
		Secure:   false,
		Path:     "/",
		HttpOnly: false,
		Value:    user.Avatar,
		Expires:  time.Now().Add(time.Hour * 24 * 12 * 30),
		SameSite: http.SameSiteDefaultMode}
	http.SetCookie(w, &userAvatarCookie)
	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
}

// GetUserFromAuthorizationHeader from jwt
func GetUserFromAuthorizationHeader(r *http.Request) *db.Td4User {

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Bearer" {
		return nil
	}

	token, err := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("td4_jwt_secret")), nil
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

type key int

const contextKeyUser = 0

// GetUserFromContext self explanatory!
func GetUserFromContext(r *http.Request) *db.Td4User {
	return r.Context().Value(key(contextKeyUser)).(*db.Td4User)
}

// WithUserInContext self explanatory!
func WithUserInContext(user *db.Td4User, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), key(contextKeyUser), user)
	return r.WithContext(ctx)
}

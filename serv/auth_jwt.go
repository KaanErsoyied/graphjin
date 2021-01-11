package serv

import (
	"context"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func jwtHandler(next http.HandlerFunc) http.HandlerFunc {
	var key interface{}

	cookie := conf.GetString("auth.cookie")

	conf.BindEnv("auth.secret", "SG_AUTH_SECRET")
	secret := conf.GetString("auth.secret")

	conf.BindEnv("auth.public_key_file", "SG_AUTH_PUBLIC_KEY_FILE")
	publicKeyFile := conf.GetString("auth.public_key_file")

	switch {
	case len(secret) != 0:
		key = []byte(secret)

	case len(publicKeyFile) != 0:
		kd, err := ioutil.ReadFile(publicKeyFile)
		if err != nil {
			panic(err)
		}

		switch conf.GetString("auth.public_key_type") {
		case "ecdsa":
			key, err = jwt.ParseECPublicKeyFromPEM(kd)

		case "rsa":
			key, err = jwt.ParseRSAPublicKeyFromPEM(kd)

		default:
			key, err = jwt.ParseECPublicKeyFromPEM(kd)

		}

		if err != nil {
			panic(err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var tok string

		if len(cookie) != 0 {
			ck, err := r.Cookie(cookie)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			tok = ck.Value
		} else {
			ah := r.Header.Get(authHeader)
			if len(ah) < 10 {
				next.ServeHTTP(w, r)
				return
			}
			tok = ah[7:]
		}

		token, err := jwt.ParseWithClaims(tok, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
			ctx := context.WithValue(r.Context(), userIDKey, claims.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		next.ServeHTTP(w, r)
	}
}

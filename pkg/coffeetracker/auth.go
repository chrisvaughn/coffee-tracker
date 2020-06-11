package coffeetracker

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type AuthContextKey int

const (
	AuthContextUserID AuthContextKey = iota
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func AuthMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := os.Getenv("AUTH0_AUD")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := os.Getenv("AUTH0_ISS")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("AUTH0_ISS") + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func GetUserID(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := ctx.Value("user").(*jwt.Token)
		if sub, ok := token.Claims.(jwt.MapClaims)["sub"]; ok {
			r = r.WithContext(context.WithValue(ctx, AuthContextUserID, sub))
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//func checkScope(scope string, tokenString string) bool {
//	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		cert, err := getPemCert(token)
//		if err != nil {
//			return nil, err
//		}
//		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
//		return result, nil
//	})
//
//	claims, ok := token.Claims.(*CustomClaims)
//
//	hasScope := false
//	if ok && token.Valid {
//		result := strings.Split(claims.Scope, " ")
//		for i := range result {
//			if result[i] == scope {
//				hasScope = true
//			}
//		}
//	}
//
//	return hasScope
//}

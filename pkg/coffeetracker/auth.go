package coffeetracker

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/chrisvaughn/coffeetracker/pkg/httputils"
	"github.com/dgrijalva/jwt-go"
)

type AuthContextKey int

const (
	AuthContextUser AuthContextKey = iota
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKey `json:"keys"`
}

type JSONWebKey struct {
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

var publicKeyCache map[string]*rsa.PublicKey

func init() {
	publicKeyCache = make(map[string]*rsa.PublicKey)
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

			publicKey, err := getPublicKey(token)
			if err != nil {
				panic(err.Error())
			}

			return publicKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getPublicKey(token *jwt.Token) (*rsa.PublicKey, error) {

	if publicKey, ok := publicKeyCache[token.Header["kid"].(string)]; ok {
		return publicKey, nil
	}

	resp, err := http.Get(os.Getenv("AUTH0_ISS") + ".well-known/jwks.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		cert := "-----BEGIN CERTIFICATE-----\n" + key.X5c[0] + "\n-----END CERTIFICATE-----"
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, err
		}
		publicKeyCache[key.Kid] = publicKey
	}

	if publicKey, ok := publicKeyCache[token.Header["kid"].(string)]; ok {
		return publicKey, nil
	}
	return nil, fmt.Errorf("jwk for kid not found")

}

func (s *Service) GetUserMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := ctx.Value("user").(*jwt.Token)
		if sub, ok := token.Claims.(jwt.MapClaims)["sub"].(string); ok {
			fmt.Printf("%s\n", sub)
			user, err := s.storage.GetOrCreateUser(ctx, sub)
			if err != nil {
				httputils.ErrorResponse(w, err.Error(), 500)
			}
			r = r.WithContext(context.WithValue(ctx, AuthContextUser, user))
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//func checkScope(scope string, tokenString string) bool {
//	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		cert, err := getPublicKey(token)
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

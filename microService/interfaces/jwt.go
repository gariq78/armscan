package interfaces

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type JWTParams struct {
	Username        string
	Password        string
	LifetimeMinutes int
	Secret          []byte
}

func (ms *MicroServiceV1) GetToken(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != ms.JWTParams.Username || password != ms.JWTParams.Password {
		http.Error(w, "Неверное имя пользователя или пароль.", http.StatusUnauthorized)
		return
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(time.Duration(ms.JWTParams.LifetimeMinutes) * time.Minute).Unix(),
	})

	tokenEncoded, err := token.SignedString(ms.JWTParams.Secret)
	if err != nil {
		ms.serverError(w, fmt.Sprintf("token.SignedString: %s", err.Error()))
		return
	}

	ms.Logger.Inf("gettoken=%s", tokenEncoded)

	ms.serverInfo(w, struct {
		Access_token string `json:"access_token"`
	}{
		Access_token: tokenEncoded,
	})
}

func (ms *MicroServiceV1) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.Logger.Inf("JWTMiddleware path=%s, Authorization=%s", r.URL.Path, r.Header.Get("Authorization"))
		token, err := request.ParseFromRequest(r, bearerRemover{"Authorization"}, ms.signingKeyFunc)
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					msg := fmt.Sprintf("invalid token: %s", err.Error())
					ms.Logger.Inf(msg)
					http.Error(w, msg, http.StatusUnauthorized)
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					msg := fmt.Sprintf("expired token: %s", err.Error())
					ms.Logger.Inf(msg)
					http.Error(w, msg, http.StatusUnauthorized)
				} else {
					msg := fmt.Sprintf("token error: %s", err.Error())
					ms.Logger.Inf(msg)
					http.Error(w, msg, http.StatusUnauthorized)
				}
			} else {
				msg := fmt.Sprintf("token: %s", err.Error())
				ms.Logger.Inf(msg)
				http.Error(w, msg, http.StatusUnauthorized)
			}

			return
		}

		if !token.Valid {
			msg := "invalid token"
			ms.Logger.Inf(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (ms *MicroServiceV1) signingKeyFunc(t *jwt.Token) (interface{}, error) {
	return ms.JWTParams.Secret, nil
}

type bearerRemover []string

func (e bearerRemover) ExtractToken(req *http.Request) (string, error) {
	for _, header := range e {
		if ah := req.Header.Get(header); ah != "" {
			return strings.TrimPrefix(strings.TrimPrefix(ah, "Bearer "), "bearer "), nil
		}
	}
	return "", request.ErrNoTokenInRequest
}

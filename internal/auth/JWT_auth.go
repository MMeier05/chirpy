package auth

import (
	"errors"
	"strings"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claim := jwt.RegisteredClaims{
			Issuer: 		"chirpy-access",
			IssuedAt: 	jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: 	jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Subject: 		userID.String(),
		}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) 
	signedJWT, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedJWT, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
	return uuid.Nil, err
	}
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
	if !found {
		return "", errors.New("No Authorization Header found")
	}
	return tokenString, nil
}

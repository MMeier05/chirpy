package auth

import (
	"testing"
	"github.com/google/uuid"
	"time"
	"log"
)


func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "Secret123"
	jwtToken, err := MakeJWT(userID, tokenSecret, time.Duration(1 * time.Hour))
	if err != nil {
		t.Errorf("Token: '%s' returns error: %v", jwtToken, err)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "Secret123"
	jwtToken, err := MakeJWT(userID, tokenSecret, time.Duration(1 * time.Hour))
	if err != nil {
		log.Fatal(err)
	}
	validatedID, err := ValidateJWT(jwtToken, tokenSecret)
	if userID != validatedID || err != nil {
		t.Errorf("Token: '%s' returns error: '%v'", jwtToken, err)
	}
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "Secret123"
	jwtToken, err := MakeJWT(userID, tokenSecret, time.Duration(0 * time.Hour))
	if err != nil {
		log.Fatal(err)
	}
	_, err = ValidateJWT(jwtToken, tokenSecret)
	if  err == nil {
		t.Errorf("Token: '%s' returns error: '%v'", jwtToken, err)
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "FirstSecret321"
	jwtToken, err := MakeJWT(userID, tokenSecret, time.Duration(0 * time.Hour))
	if err != nil {
		log.Fatal(err)
	}
	_, err = ValidateJWT(jwtToken, "SecondSecret123")
	if  err == nil {
		t.Errorf("Token: '%s' returns error: '%v'", jwtToken, err)
	}
}

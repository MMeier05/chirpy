package auth

import (
	"testing"
	"log"
)
func TestHashPassword(t *testing.T) {
	password := "Password123"
	_, err := HashPassword(password)
	if err != nil {
		t.Errorf("Password '%s' returns error %v", password, err)
	}
}

func TestCheckPasswordHealth(t *testing.T) {
	hashedPW, err := HashPassword("Password123")
	if err != nil {
		log.Fatal(err)
	}
	isSame, err := CheckPasswordHash("Password123", hashedPW)
	if isSame != true || err != nil {
		t.Errorf("want isSame == true; got %t, error: %v", isSame, err)
	}
}

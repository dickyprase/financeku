package hash

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashed == "" {
		t.Fatal("HashPassword returned empty string")
	}

	if hashed == password {
		t.Fatal("HashPassword returned plain text password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(password, hashed) {
		t.Fatal("CheckPassword should return true for correct password")
	}

	if CheckPassword("wrongpassword", hashed) {
		t.Fatal("CheckPassword should return false for wrong password")
	}
}

func TestCheckPasswordEmpty(t *testing.T) {
	hashed, _ := HashPassword("somepassword")

	if CheckPassword("", hashed) {
		t.Fatal("CheckPassword should return false for empty password")
	}
}

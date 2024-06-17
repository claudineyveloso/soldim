package auth

import (
	"bytes"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if len(hash) == 0 {
		t.Error("expected hash to be not empty")
	}

	if bytes.Equal([]byte(hash), []byte(password)) {
		t.Error("expected hash to be different from password")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	hashStr := string(hash)

	if !ComparePasswords(hashStr, []byte("password")) {
		t.Errorf("expected password to match hash")
	}
	if ComparePasswords(hashStr, []byte("notpassword")) {
		t.Errorf("expected password to not match hash")
	}
}

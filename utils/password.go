package utils

import (
	"golang.org/x/crypto/bcrypt"
)


func HashText(txt string) (*string, error) {
	hashBytes, error := bcrypt.GenerateFromPassword([]byte(txt), HASH_COST)

	if error != nil {
		return nil, error
	}
	hash := string(hashBytes)
	return &hash, nil
}

func ValidateHash(hash string, txt string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(txt)); err != nil {
		return false
	}
	return true

}

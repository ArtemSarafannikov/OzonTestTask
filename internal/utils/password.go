package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword This function maded as a variable for tests
var HashPassword = func(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash This function maded as a variable for tests
var CheckPasswordHash = func(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//func HashPassword(password string) (string, error) {
//	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
//	return string(bytes), err
//}
//
//func CheckPasswordHash(password, hash string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//	return err == nil
//}

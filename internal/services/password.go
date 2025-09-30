package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Hash_this_password(password string) string {
	hashed_pass_byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("error while generate the hased password - ", err)
		return ""
	}
	hashed_pass_string := string(hashed_pass_byte)
	return hashed_pass_string
}

func Is_password_matched(hashed_password string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

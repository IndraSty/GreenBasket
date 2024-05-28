package util

import (
	"log"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "login or password incorrect"
		check = false
	}

	return check, msg
}

func ValidatePassword(password string) string {
	var (
		errors []string
	)

	if !regexp.MustCompile(`\d`).MatchString(password) {
		errors = append(errors, "Password must contain a number")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		errors = append(errors, "Password must contain a lowercase")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		errors = append(errors, "Password must contain an uppercase")
	}
	if !regexp.MustCompile(`[!@#$%^&*()_+]`).MatchString(password) {
		errors = append(errors, "Password must contain a special character")
	}

	return strings.Join(errors, ", ")
}

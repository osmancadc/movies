package main

import (
	"errors"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/robbert229/jwt"
)

func ValidateParameters(user User) (bool, error) {
	validate := validator.New()
	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return len(password) > 9 && HasUpper(password) && HasLower(password) && HasSpecialCharacters(password)
	})

	err := validate.Var(user.Email, "required,email")
	if err != nil {
		return false, errors.New("invalid email address")
	}

	err = validate.Var(user.Password, "required,password")
	if err != nil {
		return false, errors.New("the password is too weak")
	}

	return true, nil
}

func HasUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func HasLower(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func HasSpecialCharacters(s string) bool {
	return strings.Contains(s, "!") ||
		strings.Contains(s, "]") ||
		strings.Contains(s, "?") ||
		strings.Contains(s, "#") ||
		strings.Contains(s, "@")
}

func GenerateToken(email string) (string, error) {
	if len(email) == 0 {
		return "", errors.New("invalid email")
	}

	algorithm := jwt.HmacSha256("SecretPassword")
	ttl := 20 * time.Minute

	claims := jwt.NewClaim()
	claims.Set("email", email)
	claims.Set("exp", time.Now().UTC().Add(ttl).Unix())

	token, err := algorithm.Encode(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(token string) bool {
	algorithm := jwt.HmacSha256("SecretPassword")

	if algorithm.Validate(token) != nil {
		log.Println(algorithm.Validate(token))
		return false
	}
	return true
}

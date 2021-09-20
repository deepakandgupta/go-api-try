package helpers

import (
	"net/mail"
	"unicode"
)

func CheckValidEmail(email string) error {
    _, err := mail.ParseAddress(email)
    return err
}

func CheckIfValidPassword(s string) bool {
    var (
        hasMinLen  = false
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )
    if len(s) >= 8 {
        hasMinLen = true
    }
    for _, char := range s {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
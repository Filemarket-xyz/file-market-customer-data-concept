package handler

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"
)

func addressValidation(address string) error {
	regExp, err := regexp.Compile("^0x[a-fA-F0-9]{40}$")
	if err != nil {
		return fmt.Errorf("addressValidation/regexp.Compile: %w", err)
	}
	if !regExp.MatchString(address) {
		return errors.New("invalid address")
	}
	return nil
}

func loginValidation(login string) error {
	if login == "" {
		return nil
	}
	if len(login) < 3 || len(login) > 20 {
		return errors.New("loginValidation incorrect length")
	}

	for _, r := range login {
		if !unicode.IsDigit(r) && !unicode.Is(unicode.Latin, r) {
			return errors.New("loginValidation incorrect login")
		}
	}
	return nil
}

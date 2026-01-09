package services

import (
	"regexp"
	"strings"
)

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	lower := regexp.MustCompile(`[a-z]`).MatchString
	upper := regexp.MustCompile(`[A-Z]`).MatchString
	number := regexp.MustCompile(`[0-9]`).MatchString
	special := regexp.MustCompile(`[@$!%*?&]`).MatchString
	return lower(password) && upper(password) && number(password) && special(password)
}

func ValidateClgMail(email string) bool {
	var is_valid bool
	isMatch, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if !isMatch {
		is_valid = false
	} else {
		isGmail := strings.Split(email, "@")
		if isGmail[1] != "kamarajengg.edu.in" {
			is_valid = false
		} else {
			is_valid = true
		}
	}
	return is_valid
}

// validate the name
func ValidateName(name string) bool {
	// Allows alphabets and spaces, 2 to 50 chars long
	re := regexp.MustCompile(`^[A-Za-z ]{2,50}$`)
	is_valid := re.MatchString(name)
	return is_valid
}

// validate the mobile_no with the regexp

func ValidateMobileNo(mobileNo string) bool {
	re := regexp.MustCompile(`^[6-9]\d{9}$`)
	is_valid := re.MatchString(mobileNo)
	return is_valid
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	is_valid := re.MatchString(email)
	return is_valid
}

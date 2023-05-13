package validator

import (
	"net/mail"
	"regexp"
	"unicode/utf8"

	"forum/internal/models"
)

const (
	MsgEmailExists        = "Email already exists"
	MsgNameExists         = "Name already exists"
	MsgInvalidEmail       = "Write correct email"
	MsgInvalidName        = "Write correct name"
	MsgInvalidPass        = "Password must contain letters, numbers and must be at least 6 characters."
	MsgUserNotFound       = "User not found"
	MsgPassDontMatch      = "The passwords don't match"
	MsgNotCorrectPassword = "Password is not correct"
)

func GetErrMsgs(m models.User) map[string]string {
	errmsgs := make(map[string]string)
	if !isValidEmail(m.Email) {
		errmsgs["email"] = MsgInvalidEmail
	}
	if !isValidName(m.Name) {
		errmsgs["name"] = MsgInvalidName
	}
	if !isValidPassword(m.Password) {
		errmsgs["pass"] = MsgInvalidPass
	}
	if m.Password != m.ConPassword {
		errmsgs["conpass"] = MsgPassDontMatch
	}
	return errmsgs
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidName(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 5 || length > 15 {
		return false
	}
	usernameConvention := "^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$"
	re, _ := regexp.Compile(usernameConvention)
	return re.MatchString(name)
}

func isValidPassword(pass string) bool {
	tests := []string{".{6,}", "[a-z]", "[0-9]"}
	for _, test := range tests {
		valid, _ := regexp.MatchString(test, pass)
		if !valid {
			return false
		}
	}
	return true
}

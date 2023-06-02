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
	MsgInvalidName        = "Write correct name. Username should start with an alphabet [A-Za-z] and all other characters can be alphabets, numbers or an underscore so, [A-Za-z0-9_]. The username consists of 5 to 15 characters inclusive."
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
	rxEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidName(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 5 || length > 15 {
		return false
	}
	usernameConvention := "^[A-Za-z][A-Za-z0-9_]{4,14}$"
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

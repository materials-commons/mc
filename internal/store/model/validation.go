package model

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	userDotRegexp = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
	userRegexp    = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp    = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
)

// IsEmail validates an email address. The code is from github.com/asaskevich/govalidator, however
// govalidator also attempts to validate that the host exists. That check was not something we needed.
func IsEmail(value interface{}) error {
	email, ok := value.(string)
	if !ok {
		return fmt.Errorf("email must be a string: '%#v'", value)
	}

	if len(email) < 6 || len(email) > 254 {
		return fmt.Errorf("bad email length '%s'", email)
	}
	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return fmt.Errorf("bad email format @ missing or badly placed '%s'", email)
	}
	user := email[:at]
	host := email[at+1:]
	if len(user) > 64 {
		return fmt.Errorf("email username too long %d", len(user))
	}

	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return fmt.Errorf("badly formed email address '%s'", email)
	}

	return nil
}

/*

err := validation.Validate("xyz", validation.By(checkAbc))
*/

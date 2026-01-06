package core

func IsValidPassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordLength
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSymbol := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		case (c >= 33 && c <= 47) || (c >= 58 && c <= 64) || (c >= 91 && c <= 96) || (c >= 123 && c <= 126):
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSymbol {
		return ErrWeakPassword
	}

	return nil
}

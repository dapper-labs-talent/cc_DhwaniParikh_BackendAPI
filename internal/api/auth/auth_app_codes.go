package auth

// Authentication relate app status codes
// Reserved range of 1 to 999
const (
	appCodeBindFailed         = 1
	appCodeCreateError        = 2
	appCodeAccountLoadError   = 3
	appCodeAccountUpdateError = 4
	appCodeAccountDeleteError = 5
	appCodeTokenError         = 6
)

var codeText = map[int]string{
	appCodeBindFailed:         "Unable to bind request",
	appCodeCreateError:        "Unable to create account",
	appCodeAccountLoadError:   "Unable to retrieve account information",
	appCodeAccountUpdateError: "Unable to update account",
	appCodeAccountDeleteError: "Unable to delete account or related data",
}

func appCodeText(code int) string {
	return codeText[code]
}

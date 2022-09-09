package appuser

// User app status codes
// Reserved range of 1 to 999
const (
	appCodeBindFailed         = 1
	appCodeAccountLoadError   = 2
	appCodeAccountUpdateError = 3
)

var codeText = map[int]string{
	appCodeBindFailed:         "Unable to bind request",
	appCodeAccountLoadError:   "Unable to retrieve account information",
	appCodeAccountUpdateError: "Unable to update account",
}

func appCodeText(code int) string {
	return codeText[code]
}

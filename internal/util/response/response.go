package response

import (
	"github.com/go-chi/render"
	"net/http"
)

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err error `json:"-"` // low-level runtime error

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

//--
// Utility responses and Error response payloads
//--

//--
// 200 responses
//--

// Respond with 200 Ok
func Ok(w http.ResponseWriter, r *http.Request, v interface{}) {
	respondWithStatus(w, r, http.StatusOK, v)
}

// Respond with 201 Created
func Created(w http.ResponseWriter, r *http.Request, v interface{}) {
	respondWithStatus(w, r, http.StatusCreated, v)
}

// Respond with 204 NoContent
func NoContent(w http.ResponseWriter, r *http.Request) {
	respondWithStatus(w, r, http.StatusNoContent, nil)
}

// Respond with 400 BadRequest
func BadRequest(w http.ResponseWriter, r *http.Request, e ErrResponse) {
	respondWithStatus(w, r, http.StatusBadRequest, e)
}

// Respond with 404 Not Found
func NotFound(w http.ResponseWriter, r *http.Request) {
	respondWithStatus(w, r, http.StatusNotFound, nil)
}

// Respond with 401 Unauthorized
func Unauthorized(w http.ResponseWriter, r *http.Request) {
	respondWithStatus(w, r, http.StatusUnauthorized, nil)
}

// Respond with 403 Forbidden
func Forbidden(w http.ResponseWriter, r *http.Request, e ErrResponse) {
	respondWithStatus(w, r, http.StatusForbidden, e)
}

// Respond with 405 MethodNotAllowed
func MethodNotAllowed(w http.ResponseWriter, r *http.Request, e ErrResponse) {
	respondWithStatus(w, r, http.StatusMethodNotAllowed, e)
}

// Respond with 500 InternalServerError
func InternalServerError(w http.ResponseWriter, r *http.Request, e ErrResponse) {
	respondWithStatus(w, r, http.StatusInternalServerError, e)
}

// CustomError responds with a the provided HTTP status code and ErrResponse
func CustomError(w http.ResponseWriter, r *http.Request, statusCode int, e ErrResponse) {
	respondWithStatus(w, r, statusCode, e)
}

func respondWithStatus(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	render.Status(r, status)
	render.JSON(w, r, v)
}

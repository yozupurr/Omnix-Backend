package httpx

import "net/http"

type Error struct {
	Status      int      `json:"-"`
	Code        string   `json:"errorCode"`
	Message     string   `json:"errorMessage"`
	MessageVars []string `json:"messageVars"`
	NumericCode int      `json:"numericErrorCode"`
	Service     string   `json:"originatingService"`
	Intent      string   `json:"intent"`
}

var (
	NotFound         = newError(http.StatusNotFound, "errors.com.epicgames.common.not_found", "Resource not found", 1004)
	InvalidRequest   = newError(http.StatusBadRequest, "errors.com.epicgames.common.oauth.invalid_request", "Invalid request", 1013)
	UnsupportedGrant = newError(http.StatusBadRequest, "errors.com.epicgames.common.oauth.unsupported_grant_type", "Unsupported grant type", 1016)
	InvalidClient    = newError(http.StatusUnauthorized, "errors.com.epicgames.common.oauth.invalid_client", "Invalid client credentials", 18033)
	AccountNotFound  = newError(http.StatusNotFound, "errors.com.epicgames.account.account_not_found", "Account not found", 18007)
	InvalidToken     = newError(http.StatusUnauthorized, "errors.com.epicgames.common.authentication.token_verification_failed", "Token verification failed", 1014)
	ValidationFailed = newError(http.StatusBadRequest, "errors.com.epicgames.validation.validation_failed", "Validation failed", 1040)
)

func OperationForbidden(profileID string) Error {
	err := newError(
		http.StatusForbidden,
		"errors.com.epicgames.modules.profiles.operation_forbidden",
		"Unable to find template configuration for profile "+profileID,
		12813,
	)
	err.MessageVars = []string{profileID}
	return err
}

func newError(status int, code, message string, numeric int) Error {
	return Error{
		Status:      status,
		Code:        code,
		Message:     message,
		MessageVars: []string{},
		NumericCode: numeric,
		Service:     "any",
		Intent:      "prod",
	}
}

func (e Error) Write(w http.ResponseWriter) {
	JSON(w, e.Status, e)
}

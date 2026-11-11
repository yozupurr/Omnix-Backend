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
	NotFound            = newError(http.StatusNotFound, "errors.com.epicgames.common.not_found", "Resource not found", 1004)
	InvalidRequest      = newError(http.StatusBadRequest, "errors.com.epicgames.common.oauth.invalid_request", "Invalid request", 1013)
	UnsupportedGrant    = newError(http.StatusBadRequest, "errors.com.epicgames.common.oauth.unsupported_grant_type", "Unsupported grant type", 1016)
	InvalidClient       = newError(http.StatusUnauthorized, "errors.com.epicgames.common.oauth.invalid_client", "Invalid client credentials", 18033)
	AccountNotFound     = newError(http.StatusNotFound, "errors.com.epicgames.account.account_not_found", "Account not found", 18007)
	InvalidRefreshToken = newError(http.StatusBadRequest, "errors.com.epicgames.account.auth_token.invalid_refresh_token", "Invalid refresh token", 18036)
	InvalidToken        = newError(http.StatusUnauthorized, "errors.com.epicgames.common.authentication.token_verification_failed", "Token verification failed", 1014)
	ValidationFailed    = newError(http.StatusBadRequest, "errors.com.epicgames.validation.validation_failed", "Validation failed", 1040)
	CatalogOutOfDate    = newError(http.StatusBadRequest, "errors.com.epicgames.modules.gamesubcatalog.catalog_out_of_date", "Offer is no longer available", 28001)
	PurchaseNotAllowed  = newError(http.StatusBadRequest, "errors.com.epicgames.modules.gamesubcatalog.purchase_not_allowed", "Purchase not allowed", 28004)
	NotEnoughVBucks     = newError(http.StatusBadRequest, "errors.com.epicgames.modules.gameplayutils.not_enough_mtx", "Not enough V-Bucks", 12720)
	PartyNotFound       = newError(http.StatusNotFound, "errors.com.epicgames.social.party.party_not_found", "Party not found", 51002)
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

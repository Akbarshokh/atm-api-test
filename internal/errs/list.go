package errs

import "errors"

var (
	ErrInternal      = errors.New("internal error")
	ErrValidation    = errors.New("validation error")
	ErrExternal      = errors.New("external service is died")
	ErrAuthorization = errors.New("authorization failed")
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token expired")
)

var (
	ErrProviderNotAllowed     = errors.New("provider_not_allowed")
	ErrCardBlocked            = errors.New("card_blocked")
	ErrInvalidPIN             = errors.New("invalid_pin")
	ErrPinAttemptsExceeded    = errors.New("pin_attempts_exceeded")
	ErrAmountNotMultiple10    = errors.New("amount_not_multiple_of_10")
	ErrInsufficientFunds      = errors.New("insufficient_funds")
	ErrCardNotFound           = errors.New("card_not_found")
	ErrConcurrentModification = errors.New("concurrent_modification")
)

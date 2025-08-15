package errs

import "errors"

var (
	ErrInternal           = errors.New("internal error")
	ErrValidation         = errors.New("validation error")
	ErrExternal           = errors.New("external service is died")
	ErrAuthorization      = errors.New("authorization failed")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenEmpty         = errors.New("empty token")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailAlreadyExist  = errors.New("email already exists")
	ErrNotFound           = errors.New("entity not found")
	ErrSurveyNotComplete  = errors.New("survey not complete")
	ErrDuplicateKey       = errors.New("duplicate key")
	ErrAlreadyLinked      = errors.New("user already has a partner")
	ErrSelfInvite         = errors.New("cannot accept your own invite")
	ErrInvalidInviteCode  = errors.New("invalid or expired invite code")
	ErrInvalidOldPassword = errors.New("invalid old password")
)

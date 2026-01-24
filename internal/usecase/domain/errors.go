package domain

import "errors"

var (
	ErrGetAccounts         = errors.New("error getting accounts")
	ErrGetAccount          = errors.New("error getting account")
	ErrGetTypes            = errors.New("error getting acceptable types")
	ErrCreateAccount       = errors.New("error creating account")
	ErrUpdateAccount       = errors.New("error updating account")
	ErrDeleteAccount       = errors.New("error deleting account")
	ErrDeleteAccountData   = errors.New("error deleting account data")
	ErrNotFoundAccount     = errors.New("account not found")
	ErrNilAccount          = errors.New("account is nil")
	ErrInvalidAccount      = errors.New("account data is invalid")
	ErrInvalidParams       = errors.New("request parameters are invalid")
	ErrInternal            = errors.New("internal error")
	ErrNotFoundAccountData = errors.New("account data not found")
	ErrGetAccountData      = errors.New("error getting account data")
	ErrGetDeals            = errors.New("error getting deals")
)

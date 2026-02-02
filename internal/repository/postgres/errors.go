package postgres

import "errors"

var (
	ErrGetAccounts         = errors.New("error getting accounts")
	ErrScanAccount         = errors.New("error scanning account")
	ErrNotFoundAccount     = errors.New("account not found")
	ErrNotFoundAccountData = errors.New("account data not found")
	ErrGetAccount          = errors.New("error getting account")
	ErrGetTypes            = errors.New("error getting acceptable types")
	ErrScanType            = errors.New("error scanning acceptable type")
	ErrCreateAccount       = errors.New("error creating account")
	ErrUpdateAccount       = errors.New("error updating account")
	ErrDeleteAccountData   = errors.New("error delete account data")
	ErrDeleteAccount       = errors.New("error deleting account")
	ErrNilAccount          = errors.New("account is nil")
	ErrGetAccountData      = errors.New("error getting account data")
	ErrGetDeals            = errors.New("error getting deals")
	ErrScanDeals           = errors.New("error scanning deals")
	ErrNotFoundBotLink     = errors.New("error found bot link")
)

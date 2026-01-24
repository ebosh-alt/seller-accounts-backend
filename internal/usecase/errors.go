package usecase

import "sellers-accounts-backend/internal/usecase/domain"

var (
	ErrGetAccounts         = domain.ErrGetAccounts
	ErrGetAccount          = domain.ErrGetAccount
	ErrGetTypes            = domain.ErrGetTypes
	ErrCreateAccount       = domain.ErrCreateAccount
	ErrUpdateAccount       = domain.ErrUpdateAccount
	ErrDeleteAccount       = domain.ErrDeleteAccount
	ErrNotFoundAccount     = domain.ErrNotFoundAccount
	ErrNotFoundAccountData = domain.ErrNotFoundAccountData
	ErrNilAccount          = domain.ErrNilAccount
	ErrInvalidAccount      = domain.ErrInvalidAccount
	ErrInvalidParams       = domain.ErrInvalidParams
	ErrInternal            = domain.ErrInternal
	ErrGetAccountData      = domain.ErrGetAccountData
	ErrDeleteAccountData   = domain.ErrDeleteAccountData
	ErrGetDeals            = domain.ErrGetDeals
)

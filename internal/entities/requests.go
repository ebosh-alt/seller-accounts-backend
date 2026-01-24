package entities

type RequestAllAccounts struct {
	Limit int `form:"limit" json:"limit"`
	Page  int `form:"page" json:"page"`
}

type RequestAcceptableTypeAccounts struct{}

type RequestAccountByID struct {
	ID string `json:"id" uri:"id"`
}

type RequestCreateAccount struct {
	UID         string  `json:"id"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Data        []Data  `json:"data"`
	Accepted    bool    `json:"accepted"`
	ViewType    bool    `json:"view_type"`
	Name        string  `json:"name"`
}

type RequestCreateAccounts struct {
	Accounts []RequestCreateAccount `json:"accounts"`
}

func (a *RequestCreateAccount) ToAccount() *Account {
	return &Account{
		UID:         a.UID,
		Price:       a.Price,
		Description: a.Description,
		Data:        a.Data,
		Accepted:    a.Accepted,
		ViewType:    a.ViewType,
		Name:        a.Name,
		Category:    Category{ID: a.CategoryID},
	}
}

type RequestUpdateAccount struct {
	UID         string  `json:"id"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Data        []Data  `json:"data"`
	Accepted    bool    `json:"accepted"`
	ViewType    bool    `json:"view_type"`
	Name        string  `json:"name"`
}

func (a *RequestUpdateAccount) ToAccount() *Account {
	return &Account{
		UID:         a.UID,
		Price:       a.Price,
		Description: a.Description,
		Data:        a.Data,
		Accepted:    a.Accepted,
		ViewType:    a.ViewType,
		Name:        a.Name,
	}
}

type RequestDeactivateAccountsByName struct {
	Name string `json:"name"`
}

type RequestDeleteAccount struct {
	ID string `json:"id"`
}

// RequestGetAccountData represents a request to get account data by ID.
type RequestGetAccountData struct {
	ID int `json:"id"`
}

// RequestDeleteAccountData represents a request to delete account data by ID.
type RequestDeleteAccountData struct {
	ID int `json:"id"`
}

type RequestGetDeals struct {
	Limit int `form:"limit" json:"limit"`
	Page  int `form:"page" json:"page"`
}

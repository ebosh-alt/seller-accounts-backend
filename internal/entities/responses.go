package entities

type ResponseAllAccounts struct {
	Accounts *[]Account `json:"accounts"`
	Total    int64      `json:"total"`
}

type ResponseAcceptableTypeAccounts struct {
	Types []TypeAccount `json:"types"`
}

type ResponseAccountByID struct {
	Account *Account `json:"account"`
}

type AccountCreateFailure struct {
	Index   int    `json:"index"`
	ID      string `json:"id,omitempty"`
	Message string `json:"message"`
}

type ResponseCreateAccounts struct {
	Accounts []*Account             `json:"accounts"`
	Failed   []AccountCreateFailure `json:"failed"`
	Total    int                    `json:"total"`
}

type ResponseUpdateAccount struct {
	Account *UpdateAccount `json:"account"`
}

type ResponseDeactivateAccounts struct {
	Updated int `json:"updated"`
}

type ResponseGetAccountData struct {
	Data Data `json:"data"`
}

type ResponseDeleteAccountData struct {
	Status bool `json:"status"`
}

type ResponseGetDeals struct {
	Deals *[]Deal `json:"deals"`
	Total int64   `json:"total"`
}

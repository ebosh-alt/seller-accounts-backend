package entities

import (
	"fmt"
	"time"
)

type Category struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (c Category) IsZero() bool {
	return c.ID == 0
}

type Deal struct {
	ID            *int       `json:"id,omitempty"`
	AccountID     *string    `json:"account_id,omitempty"`
	BuyerID       *int       `json:"buyer_id,omitempty"`
	Commission    *float64   `json:"commission,omitempty"`
	Data          *Data      `json:"data,omitempty"`
	Price         *float64   `json:"price,omitempty"`
	Wallet        *string    `json:"wallet,omitempty"`
	IsGuarantor   *bool      `json:"guarantor,omitempty"`
	PaymentStatus *bool      `json:"payment_status,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}

func (d *Deal) IsZero() bool {
	return d.ID == nil
}

type Account struct {
	UID         string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Link        string    `json:"link,omitempty"`
	Category    Category  `json:"category,omitzero"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Data        []Data    `json:"data"`
	Accepted    bool      `json:"accepted"`
	ViewType    bool      `json:"view_type"`
	Name        string    `json:"name"`
	Deal        Deal      `json:"deal,omitzero"`
}

func (a *Account) InitLink(botLink string) {
	if a != nil {
		a.Link = fmt.Sprintf("%s/%s", botLink, a.UID)
	}
}

type AllAccounts struct {
	Accounts   *[]Account `json:"accounts"`
	TotalCount int64      `json:"total_count"`
}

type UpdateAccount struct {
	UID         string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Category    Category  `json:"category,omitzero"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Data        []Data    `json:"data"`
	Accepted    bool      `json:"accepted"`
	ViewType    bool      `json:"view_type"`
	Name        string    `json:"name"`
}

func (a *Account) ToUpdate() *UpdateAccount {
	return &UpdateAccount{
		UID:         a.UID,
		CreatedAt:   a.CreatedAt,
		Category:    a.Category,
		Price:       a.Price,
		Description: a.Description,
		Data:        a.Data,
		Accepted:    a.Accepted,
		ViewType:    a.ViewType,
		Name:        a.Name,
	}
}

type TypeAccount struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AcceptableTypesAccounts []TypeAccount

type Data struct {
	ID        int    `json:"id,omitempty"`
	AccountID string `json:"account_id"`
	DealID    *int   `json:"deal_id,omitempty"`
	IsPayment bool   `json:"is_payment"`
	Value     string `json:"value"`
}

type Accounts []*Account

func (a Accounts) InitLinks(baseURL string) {
	for i := range a {
		a[i].InitLink(baseURL)
	}
}

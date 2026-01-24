package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"sellers-accounts-backend/internal/entities"
)

const (
	AllAccountsQuery = `SELECT COUNT(*) OVER() AS total, accounts.id, accounts.created_at, accounts.name, accounts.price, accounts.description, data_list.data, accounts.accepted, accounts.view_type, accounts.category_id, c.name, d.id, d.buyer_id, d.price, d.wallet, d.payment_status, d.created_at, d.guarantor
FROM accounts
LEFT JOIN categories c ON c.id = accounts.category_id
LEFT JOIN LATERAL (
	SELECT json_agg(json_build_object('id', id, 'account_id', account_id, 'is_payment', is_payment, 'value', value) ORDER BY id) AS data
	FROM accounts_data
	WHERE account_id = accounts.id
) data_list ON true
LEFT JOIN LATERAL (
	SELECT id
	FROM accounts_data
	WHERE account_id = accounts.id
	ORDER BY id DESC
	LIMIT 1
) last_data ON true
LEFT JOIN deals d ON d.data_id = last_data.id
ORDER BY accounts.view_type DESC, accounts.created_at DESC
LIMIT $1 OFFSET $2`
	AcceptableTypesAccountsQuery = `SELECT id, name FROM categories;`
	AccountByIDQuery             = `SELECT accounts.id, accounts.created_at, accounts.name, accounts.price, accounts.description, data_list.data, accounts.accepted, accounts.view_type, accounts.category_id, c.name, d.id, d.buyer_id, d.price, d.wallet, d.payment_status, d.created_at, d.guarantor
FROM accounts
LEFT JOIN categories c ON c.id = accounts.category_id
LEFT JOIN LATERAL (
	SELECT json_agg(json_build_object('id', id, 'account_id', account_id, 'is_payment', is_payment, 'value', value) ORDER BY id) AS data
	FROM accounts_data
	WHERE account_id = accounts.id
) data_list ON true
LEFT JOIN LATERAL (
	SELECT id
	FROM accounts_data
	WHERE account_id = accounts.id
	ORDER BY id DESC
	LIMIT 1
) last_data ON true
LEFT JOIN deals d ON d.data_id = last_data.id
WHERE accounts.id = $1;`
	CreateAccountQuery            = `WITH inserted AS (INSERT INTO accounts (category_id, name, price, description, accepted, view_type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, category_id, created_at) SELECT inserted.id, c.id, c.name, inserted.created_at FROM inserted JOIN categories c ON c.id = inserted.category_id;`
	UpdateAccountQuery            = `UPDATE accounts SET name = $1, price = $2, description = $3, accepted = $4, view_type = $5 WHERE id = $6;`
	UpdateAccountDataByIDQuery    = `UPDATE accounts_data SET is_payment = $1, value = $2 WHERE id = $3 AND account_id = $4;`
	InsertAccountDataQuery        = `INSERT INTO accounts_data (account_id, is_payment, value) VALUES ($1, $2, $3) RETURNING id;`
	DeactivateAccountsByNameQuery = `UPDATE accounts SET view_type = false WHERE name = $1;`
	DeleteDealsByAccountQuery     = `DELETE FROM deals WHERE data_id IN (SELECT id FROM accounts_data WHERE account_id = $1);`
	DeleteAccountDataByAccount    = `DELETE FROM accounts_data WHERE account_id = $1;`
	DeleteAccountQuery            = `DELETE FROM accounts WHERE id = $1;`
)

func (p *Postgres) AllAccounts(ctx context.Context, limit int, page int) (*[]entities.Account, int64, error) {
	offset := offSet(limit, page)

	rows, err := p.db.Query(ctx, AllAccountsQuery, limit, offset)
	if err != nil {
		p.log.Errorw("AllAccounts - p.db.Query", "error", err)
		return nil, 0, ErrGetAccounts
	}
	defer rows.Close()

	accounts := make([]entities.Account, 0, 61)
	var totalCount int64
	for rows.Next() {
		var a entities.Account
		var dataJSON []byte
		err := rows.Scan(
			&totalCount,
			&a.UID,
			&a.CreatedAt,
			&a.Name,
			&a.Price,
			&a.Description,
			&dataJSON,
			&a.Accepted,
			&a.ViewType,
			&a.Category.ID,
			&a.Category.Name,
			&a.Deal.ID,
			&a.Deal.BuyerID,
			&a.Deal.Price,
			&a.Deal.Wallet,
			&a.Deal.PaymentStatus,
			&a.Deal.CreatedAt,
			&a.Deal.IsGuarantor,
		)
		if err != nil {
			p.log.Errorw("AllAccounts - rows.Scan", "error", err)
			return nil, 0, ErrScanAccount
		}
		if len(dataJSON) > 0 && string(dataJSON) != "null" {
			if err := json.Unmarshal(dataJSON, &a.Data); err != nil {
				p.log.Errorw("AllAccounts - json.Unmarshal", "error", err)
				return nil, 0, ErrScanAccount
			}
		}
		accounts = append(accounts, a)
	}
	if err := rows.Err(); err != nil {
		p.log.Errorw("AllAccounts - rows.Err", "error", err)
		return nil, 0, ErrScanAccount
	}
	return &accounts, totalCount, nil
}

func (p *Postgres) AcceptableTypesAccounts(ctx context.Context) (entities.AcceptableTypesAccounts, error) {
	rows, err := p.db.Query(ctx, AcceptableTypesAccountsQuery)
	if err != nil {
		p.log.Errorw("AcceptableTypesAccounts - p.db.QueryRow", "error", err)
		return nil, ErrGetTypes
	}

	defer rows.Close()

	types := make([]entities.TypeAccount, 0, 61)
	for rows.Next() {
		var t entities.TypeAccount
		err := rows.Scan(
			&t.ID,
			&t.Name,
		)
		if err != nil {
			p.log.Errorw("AcceptableTypesAccounts - rows.Scan", "error", err)
			continue
		}
		types = append(types, t)
	}

	return types, nil
}

func (p *Postgres) Account(ctx context.Context, uid string) (*entities.Account, error) {
	a := entities.Account{}
	var dataJSON []byte
	err := p.db.QueryRow(ctx, AccountByIDQuery, uid).Scan(
		&a.UID,
		&a.CreatedAt,
		&a.Name,
		&a.Price,
		&a.Description,
		&dataJSON,
		&a.Accepted,
		&a.ViewType,
		&a.Category.ID,
		&a.Category.Name,
		&a.Deal.ID,
		&a.Deal.BuyerID,
		&a.Deal.Price,
		&a.Deal.Wallet,
		&a.Deal.PaymentStatus,
		&a.Deal.CreatedAt,
		&a.Deal.IsGuarantor,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Infow("account not found", "url", uid)
			return nil, ErrNotFoundAccount
		}
		p.log.Errorw("failed to get url", "error", err)
		return nil, ErrGetAccount
	}
	if len(dataJSON) > 0 && string(dataJSON) != "null" {
		if err := json.Unmarshal(dataJSON, &a.Data); err != nil {
			p.log.Errorw("Account - json.Unmarshal", "error", err)
			return nil, ErrGetAccount
		}
	}
	return &a, nil
}

func (p *Postgres) CreateAccount(ctx context.Context, a *entities.Account) (*entities.Account, error) {
	if a == nil {
		return nil, ErrNilAccount
	}
	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.log.Errorw("CreateAccount - p.db.Begin", "error", err)
		return nil, ErrCreateAccount
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	err = tx.QueryRow(
		ctx,
		CreateAccountQuery,
		a.Category.ID,
		a.Name,
		a.Price,
		a.Description,
		a.Accepted,
		a.ViewType,
	).Scan(&a.UID, &a.Category.ID, &a.Category.Name, &a.CreatedAt)

	if err != nil {
		p.log.Errorw("CreateAccount - tx.QueryRow", "error", err)
		_ = tx.Rollback(ctx)
		return nil, ErrCreateAccount
	}

	for i := range a.Data {
		data := &a.Data[i]
		p.log.Infow("CreateAccount - data", "data", data)
		err = tx.QueryRow(ctx, InsertAccountDataQuery, a.UID, data.IsPayment, data.Value).Scan(&data.ID)
		if err != nil {
			p.log.Errorw("CreateAccount - tx.QueryRow data", "error", err)
			_ = tx.Rollback(ctx)
			return nil, ErrCreateAccount
		}
	}

	if err = tx.Commit(ctx); err != nil {
		p.log.Errorw("CreateAccount - tx.Commit", "error", err)
		return nil, ErrCreateAccount
	}
	return a, nil
}

func (p *Postgres) UpdateAccount(ctx context.Context, a *entities.Account) error {
	if a == nil {
		return ErrNilAccount
	}
	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.log.Errorw("UpdateAccount - p.db.Begin", "error", err)
		return ErrUpdateAccount
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	commandTag, err := tx.Exec(
		ctx,
		UpdateAccountQuery,
		a.Name,
		a.Price,
		a.Description,
		a.Accepted,
		a.ViewType,
		a.UID,
	)
	if err != nil {
		p.log.Errorw("UpdateAccount - tx.Exec", "error", err)
		_ = tx.Rollback(ctx)
		return ErrUpdateAccount
	}
	if commandTag.RowsAffected() == 0 {
		p.log.Infow("account not found", "uid", a.UID)
		_ = tx.Rollback(ctx)
		return ErrNotFoundAccount
	}

	for i := range a.Data {
		data := &a.Data[i]
		if data.ID > 0 {
			dataCommandTag, err := tx.Exec(ctx, UpdateAccountDataByIDQuery, data.IsPayment, data.Value, data.ID, a.UID)
			if err != nil {
				p.log.Errorw("UpdateAccount - tx.Exec data", "error", err)
				_ = tx.Rollback(ctx)
				return ErrUpdateAccount
			}
			if dataCommandTag.RowsAffected() > 0 {
				continue
			}
		}

		err = tx.QueryRow(ctx, InsertAccountDataQuery, a.UID, data.IsPayment, data.Value).Scan(&data.ID)
		if err != nil {
			p.log.Errorw("UpdateAccount - tx.QueryRow insert data", "error", err)
			_ = tx.Rollback(ctx)
			return ErrUpdateAccount
		}
	}

	if err = tx.Commit(ctx); err != nil {
		p.log.Errorw("UpdateAccount - tx.Commit", "error", err)
		return ErrUpdateAccount
	}
	return nil
}

func (p *Postgres) DeactivateAccountsByName(ctx context.Context, name string) (int, error) {
	commandTag, err := p.db.Exec(ctx, DeactivateAccountsByNameQuery, name)
	if err != nil {
		p.log.Errorw("DeactivateAccountsByName - p.db.Exec", "error", err)
		return 0, ErrUpdateAccount
	}
	if commandTag.RowsAffected() == 0 {
		return 0, ErrNotFoundAccount
	}
	return int(commandTag.RowsAffected()), nil
}

func (p *Postgres) DeleteAccount(ctx context.Context, uid string) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.log.Errorw("DeleteAccount - p.db.Begin", "error", err)
		return ErrDeleteAccount
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, DeleteDealsByAccountQuery, uid)
	if err != nil {
		p.log.Errorw("DeleteAccount - tx.Exec deals", "error", err)
		_ = tx.Rollback(ctx)
		return ErrDeleteAccount
	}

	_, err = tx.Exec(ctx, DeleteAccountDataByAccount, uid)
	if err != nil {
		p.log.Errorw("DeleteAccount - tx.Exec data", "error", err)
		_ = tx.Rollback(ctx)
		return ErrDeleteAccount
	}

	commandTag, err := tx.Exec(ctx, DeleteAccountQuery, uid)
	if err != nil {
		p.log.Errorw("DeleteAccount - tx.Exec", "error", err)
		_ = tx.Rollback(ctx)
		return ErrDeleteAccount
	}
	if commandTag.RowsAffected() == 0 {
		p.log.Infow("account not found", "uid", uid)
		_ = tx.Rollback(ctx)
		return ErrNotFoundAccount
	}
	if err = tx.Commit(ctx); err != nil {
		p.log.Errorw("DeleteAccount - tx.Commit", "error", err)
		return ErrDeleteAccount
	}
	return nil
}

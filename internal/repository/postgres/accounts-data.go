package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"sellers-accounts-backend/internal/entities"
)

const (
	sqlAccountDataByID   = `SELECT account_id, deal_id, is_payment, value FROM accounts_data where id = $1`
	sqlDeleteAccountData = "DELETE FROM accounts_data WHERE id = $1"
)

func (p *Postgres) AccountData(ctx context.Context, id int) (*entities.Data, error) {
	var data entities.Data
	data.ID = id
	err := p.db.QueryRow(ctx, sqlAccountDataByID, id).Scan(
		&data.AccountID,
		&data.DealID,
		&data.IsPayment,
		&data.Value,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Infow("account data not found", "id", id)
			return nil, ErrNotFoundAccountData
		}
		p.log.Errorw("failed to get account data", "error", err)
		return nil, ErrGetAccountData
	}
	return &data, nil
}

func (p *Postgres) DeleteAccountData(ctx context.Context, id int) error {
	commandTag, err := p.db.Exec(ctx, sqlDeleteAccountData, id)
	if err != nil {
		p.log.Errorw("DeleteAccountData - p.db.Exec", "error", err)
		return ErrDeleteAccountData
	}
	if commandTag.RowsAffected() == 0 {
		return ErrNotFoundAccountData
	}
	return nil
}

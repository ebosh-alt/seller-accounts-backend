package postgres

import (
	"context"

	"sellers-accounts-backend/internal/entities"
)

const (
	AllDealsQuery = `SELECT COUNT(*) OVER() AS total,
       d.id,
       ad.account_id,
       d.buyer_id,
       d.data_id,
       d.price,
       d.wallet,
       d.guarantor,
       d.payment_status,
       d.created_at,
       ad.is_payment,
       ad.value
FROM deals AS d
JOIN accounts_data ad ON ad.id = d.data_id
ORDER BY d.created_at DESC
LIMIT $1 OFFSET $2`
)

func (p *Postgres) Deals(ctx context.Context, limit int, page int) (*[]entities.Deal, int64, error) {
	offset := offSet(limit, page)

	rows, err := p.db.Query(ctx, AllDealsQuery, limit, offset)
	if err != nil {
		p.log.Errorw("Deals - p.db.Query", "error", err)
		return nil, 0, ErrGetDeals
	}
	defer rows.Close()

	deals := make([]entities.Deal, 0, 61)
	var totalCount int64

	for rows.Next() {
		var deal entities.Deal
		var dataID int
		var isPayment bool
		var value string
		err := rows.Scan(
			&totalCount,
			&deal.ID,
			&deal.AccountID,
			&deal.BuyerID,
			&dataID,
			&deal.Price,
			&deal.Wallet,
			&deal.IsGuarantor,
			&deal.PaymentStatus,
			&deal.CreatedAt,
			&isPayment,
			&value,
		)
		if err != nil {
			p.log.Errorw("Deals - rows.Scan", "error", err)
			return nil, 0, ErrGetDeals
		}
		accountID := ""
		if deal.AccountID != nil {
			accountID = *deal.AccountID
		}
		deal.Data = &entities.Data{
			ID:        dataID,
			AccountID: accountID,
			IsPayment: isPayment,
			Value:     value,
		}
		deals = append(deals, deal)
	}
	if err := rows.Err(); err != nil {
		p.log.Errorw("Deals - rows.Err", "error", err)
		return nil, 0, ErrScanDeals
	}
	return &deals, totalCount, nil

}

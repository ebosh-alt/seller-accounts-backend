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
       ad.id,
       d.price,
       d.commission,
       d.wallet,
       d.payment_status,
       d.guarantor,
       d.created_at,
       ad.is_payment,
       ad.value
FROM deals AS d
LEFT JOIN accounts_data ad ON ad.deal_id = d.id
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
		var dataID *int
		var isPayment *bool
		var value *string
		err := rows.Scan(
			&totalCount,
			&deal.ID,
			&deal.AccountID,
			&deal.BuyerID,
			&dataID,
			&deal.Price,
			&deal.Commission,
			&deal.Wallet,
			&deal.PaymentStatus,
			&deal.IsGuarantor,
			&deal.CreatedAt,
			&isPayment,
			&value,
		)
		if err != nil {
			p.log.Errorw("Deals - rows.Scan", "error", err)
			return nil, 0, ErrGetDeals
		}
		if dataID != nil || value != nil || isPayment != nil {
			accountID := ""
			if deal.AccountID != nil {
				accountID = *deal.AccountID
			}
			data := &entities.Data{
				AccountID: accountID,
			}
			if dataID != nil {
				data.ID = *dataID
			}
			if deal.ID != nil {
				data.DealID = deal.ID
			}
			if isPayment != nil {
				data.IsPayment = *isPayment
			}
			if value != nil {
				data.Value = *value
			}
			deal.Data = data
		}
		deals = append(deals, deal)
	}
	if err := rows.Err(); err != nil {
		p.log.Errorw("Deals - rows.Err", "error", err)
		return nil, 0, ErrScanDeals
	}
	return &deals, totalCount, nil

}

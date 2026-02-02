package postgres

import "context"

const (
	SelectBotLinkQuery = `SELECT bot_link FROM shops LIMIT 1;`
)

func (p *Postgres) BotLink(ctx context.Context) (string, error) {
	var link string
	err := p.db.QueryRow(ctx, SelectBotLinkQuery).Scan(&link)
	if err != nil {
		return "", ErrNotFoundBotLink
	}
	return link, nil
}

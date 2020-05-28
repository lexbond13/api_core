package db

import (
	"context"
	"github.com/go-pg/pg/v9"
)

type queryDebug struct {
}

func (h *queryDebug) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	_, err := event.FormattedQuery()
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (h *queryDebug) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	return nil
}

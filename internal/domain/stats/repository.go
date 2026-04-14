package stats

import "context"

type Repository interface {
	Save(ctx context.Context, stats *Stats) (*Stats, error)
	GetLatest(ctx context.Context) (*Stats, error)
}

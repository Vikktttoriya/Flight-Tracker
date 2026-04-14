package flight

import "context"

type Repository interface {
	Create(ctx context.Context, flight *Flight) (*Flight, error)
	GetByID(ctx context.Context, id int64) (*Flight, error)
	List(ctx context.Context, offset, limit int) ([]*Flight, error)
	Update(ctx context.Context, flight *Flight) (*Flight, error)
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, id int64) error
}

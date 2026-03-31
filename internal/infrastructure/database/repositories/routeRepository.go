package repositories

import (
	"Api/internal/domain"
	"Api/internal/infrastructure/database"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IRouteRepository interface {
	Add(ctx context.Context, route *domain.Route) error
	AddAll(ctx context.Context, routes []*domain.Route) error
	FindById(ctx context.Context, id uuid.UUID) (*domain.Route, error)
	DeleteById(ctx context.Context, id uuid.UUID) (*domain.Route, error)
}

type RouteRepository struct {
	connPool *pgxpool.Pool
}

func NewRoutesRepository(dbConnPool *pgxpool.Pool) *RouteRepository {
	return &RouteRepository{
		connPool: dbConnPool,
	}
}

func (r *RouteRepository) Add(ctx context.Context, route *domain.Route) error {
	_, err := r.connPool.Exec(
		ctx,
		database.InsertRoute,
		route.Id, route.Start.Lng, route.Start.Lat,
		route.Finish.Lng, route.Finish.Lat)
	return err
}

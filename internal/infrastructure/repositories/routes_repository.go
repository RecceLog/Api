package repositories

import (
	"Api/internal/domain"
	"Api/internal/infrastructure/database"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// RoutesRepository handles all database operations for routes.
// It never opens transactions — it only executes queries through
// the provided Querier, which can be a pool or an active transaction.
type RoutesRepository interface {
	Insert(ctx context.Context, q database.Querier, r *domain.Route) error
	FindById(ctx context.Context, q database.Querier, id uuid.UUID) (domain.Route, error)
	FindAll(ctx context.Context, q database.Querier) ([]domain.Route, error)
	FindInRange(ctx context.Context, q database.Querier, point domain.Coordinate, rangeM int) ([]domain.Route, error)
	Delete(ctx context.Context, q database.Querier, id uuid.UUID) error
}

type routesRepository struct{}

func NewRoutesRepository() RoutesRepository {
	return &routesRepository{}
}

func (r *routesRepository) Insert(ctx context.Context, q database.Querier, route *domain.Route) error {

	// Insert route row
	_, err := q.Exec(
		ctx,
		database.InsertRoute,
		route.Id, route.Start.Lng, route.Start.Lat, route.Finish.Lng, route.Finish.Lat,
	)
	if err != nil {
		return err
	}

	// Insert waypoints
	for i, stop := range route.Waypoints {
		_, err = q.Exec(
			ctx,
			database.InsertRouteWaypoint,
			route.Id, stop.Position.Lng, stop.Position.Lat, i+1,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *routesRepository) FindById(ctx context.Context, q database.Querier, id uuid.UUID) (domain.Route, error) {

	rows, err := q.Query(ctx, database.SelectRouteById, id)
	if err != nil {
		return domain.Route{}, err
	}
	defer rows.Close()

	var route domain.Route
	var found bool
	var order uint

	for rows.Next() {
		var stopLat, stopLng *float64
		order++

		if !found {
			err = rows.Scan(
				&route.Id,
				&route.Start.Lat, &route.Start.Lng,
				&route.Finish.Lat, &route.Finish.Lng,
				&stopLat, &stopLng,
			)
			if err != nil {
				return domain.Route{}, err
			}
			found = true
			route.Waypoints = []domain.Waypoint{}
		} else {
			// Route fields are repeated in every row (LEFT JOIN),
			// scan them into throwaway variables
			var dummyId uuid.UUID
			var dummyF1, dummyF2, dummyF3, dummyF4 float64
			err = rows.Scan(
				&dummyId,
				&dummyF1, &dummyF2,
				&dummyF3, &dummyF4,
				&stopLat, &stopLng,
			)
			if err != nil {
				return domain.Route{}, err
			}
		}

		if stopLat != nil && stopLng != nil {
			route.Waypoints = append(route.Waypoints, domain.Waypoint{
				Position: domain.Coordinate{Lat: *stopLat, Lng: *stopLng},
				Order:    order,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return domain.Route{}, err
	}

	if !found {
		return domain.Route{}, pgx.ErrNoRows
	}

	return route, nil
}

func (r *routesRepository) FindAll(ctx context.Context, q database.Querier) ([]domain.Route, error) {

	rows, err := q.Query(ctx, database.SelectAllRoutes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRoutesRows(rows)
}

func (r *routesRepository) FindInRange(ctx context.Context, q database.Querier, point domain.Coordinate, rangeM int) ([]domain.Route, error) {

	rows, err := q.Query(ctx, database.SelectRoutesInRange, point.Lng, point.Lat, rangeM)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRoutesRows(rows)
}

func (r *routesRepository) Delete(ctx context.Context, q database.Querier, id uuid.UUID) error {

	_, err := q.Exec(ctx, database.DeleteRouteById, id)
	return err
}

// scanRoutesRows maps a result set (from a LEFT JOIN with waypoints)
// into a deduplicated list of routes with their waypoints filled in.
func scanRoutesRows(rows pgx.Rows) ([]domain.Route, error) {

	routesMap := make(map[uuid.UUID]*domain.Route)
	var orderedIds []uuid.UUID
	waypointOrder := make(map[uuid.UUID]uint)

	for rows.Next() {
		var route domain.Route
		var stopLat, stopLng *float64

		err := rows.Scan(
			&route.Id,
			&route.Start.Lat, &route.Start.Lng,
			&route.Finish.Lat, &route.Finish.Lng,
			&stopLat, &stopLng,
		)
		if err != nil {
			return nil, err
		}

		existing, seen := routesMap[route.Id]
		if !seen {
			route.Waypoints = []domain.Waypoint{}
			routesMap[route.Id] = &route
			orderedIds = append(orderedIds, route.Id)
			existing = routesMap[route.Id]
		}

		if stopLat != nil && stopLng != nil {
			waypointOrder[route.Id]++
			existing.Waypoints = append(existing.Waypoints, domain.Waypoint{
				Position: domain.Coordinate{Lat: *stopLat, Lng: *stopLng},
				Order:    waypointOrder[route.Id],
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Preserve insertion order
	result := make([]domain.Route, 0, len(orderedIds))
	for _, id := range orderedIds {
		result = append(result, *routesMap[id])
	}

	return result, nil
}

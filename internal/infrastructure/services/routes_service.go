package services

import (
	"Api/internal/domain"
	"Api/internal/infrastructure/database"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoutesService interface {
	Insert(ctx context.Context, r *domain.Route) error
	FindInRange(ctx context.Context, startingPoint domain.Coordinate, _range int) ([]domain.Route, error)
	FindById(cxt context.Context, id uuid.UUID) (domain.Route, error)
	FindAll(ctx context.Context) ([]domain.Route, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type routesService struct {
	connPool *pgxpool.Pool
}

func NewRoutesService(conn *pgxpool.Pool) RoutesService {
	return &routesService {
		connPool: conn,
	}
}

func (s *routesService) Insert(ctx context.Context, r *domain.Route) error {

	// Generate uuid
	routeId, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to create uuid locally", "error", err.Error())
		return err
	}
	r.Id = routeId

	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Prepare query to insert route
	_, err = tx.Exec(
		ctx,
		database.InsertRoute,
		routeId, r.Start.Lng, r.Start.Lat, r.Finish.Lng, r.Finish.Lat,
	)
	if err != nil {
		return err
	}

	// Prepare query to insert route waypoints
	for i, stop := range r.Waypoints {
		_, err = tx.Exec(
			ctx,
			database.InsertRouteWaypoint,
			routeId, stop.Position.Lng, stop.Position.Lat, i+1,
		)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit(ctx)
}

func (s *routesService) FindInRange(ctx context.Context, startingPoint domain.Coordinate, _range int) ([]domain.Route, error) {

	// Get routes rows from db
	rows, err := s.connPool.Query(
		ctx,
		database.SelectRoutesInRange,
		startingPoint.Lng, startingPoint.Lat, _range,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map rows to route list
	routesList, err := scanRoutesRows(rows)

	return routesList, err
}

func (s *routesService) FindById(ctx context.Context, id uuid.UUID) (domain.Route, error) {

	rows, err := s.connPool.Query(
		ctx,
		database.SelectRouteById,
		id,
	)
	if err != nil {
		return domain.Route{}, err
	}
	defer rows.Close()

	var route domain.Route
	var foundRoute bool

	var i uint
	for rows.Next() {
		var stopLat, stopLng *float64
		i++

		if !foundRoute {
			// First row - scan the route data
			err := rows.Scan(
				&route.Id,
				&route.Start.Lat, &route.Start.Lng,
				&route.Finish.Lat, &route.Finish.Lng,
				&stopLat, &stopLng,
			)
			if err != nil {
				return domain.Route{}, err
			}
			foundRoute = true
			route.Waypoints = []domain.Waypoint{}

			// Add first stop if exists
			if stopLat != nil && stopLng != nil {
				route.Waypoints = append(route.Waypoints, domain.Waypoint{
					Position: domain.Coordinate{Lat: *stopLat, Lng: *stopLng},
					Order: i,
				})
			}
		} else {
			// Subsequent rows - only scan stop data (skip route fields)
			var dummyId uuid.UUID
			var dummyStartLat, dummyStartLng, dummyFinishLat, dummyFinishLng float64

			err := rows.Scan(
				&dummyId,
				&dummyStartLat, &dummyStartLng,
				&dummyFinishLat, &dummyFinishLng,
				&stopLat, &stopLng,
			)
			if err != nil {
				return domain.Route{}, err
			}

			// Add stop if exists
			if stopLat != nil && stopLng != nil {
				route.Waypoints = append(route.Waypoints, domain.Waypoint{
					Position: domain.Coordinate{Lat: *stopLat, Lng: *stopLng},
					Order: i,
				})
			}
		}
	}

	if err := rows.Err(); err != nil {
		return domain.Route{}, err
	}

	if !foundRoute {
		return domain.Route{}, pgx.ErrNoRows
	}

	return route, nil
}

func (s *routesService) FindAll(ctx context.Context) ([]domain.Route, error) {

	rows, err := s.connPool.Query(
		ctx,
		database.SelectAllRoutes,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routesList, err := scanRoutesRows(rows)

	return routesList, err
}

func (s *routesService) Delete(ctx context.Context, id uuid.UUID) error {

	_, err := s.connPool.Exec(
		ctx,
		database.DeleteRouteById,
		id,
	)
	return err
}

func scanRoutesRows(rows pgx.Rows) ([]domain.Route, error) {

	routesMap := make(map[uuid.UUID]*domain.Route)
	var routesList []domain.Route

	var i uint
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

		if existingRoute, exists := routesMap[route.Id]; exists {
			i++
			// Add stop to existing route
			if stopLat != nil && stopLng != nil {
				existingRoute.Waypoints = append(existingRoute.Waypoints, domain.Waypoint{
					Position: domain.Coordinate{Lat: *stopLat, Lng: *stopLng},
					Order: i,
				})
			}
		} else {
			// New route
			route.Waypoints = []domain.Waypoint{}
			if stopLat != nil && stopLng != nil {
				route.Waypoints = append(route.Waypoints, domain.Waypoint{
					Position: domain.Coordinate{Lat: *stopLat, Lng: *stopLng},
					Order: i,
				})
			}
			routesMap[route.Id] = &route
			routesList = append(routesList, route)
		}
	}

	// Update routes in list with filled stops
	for i := range routesList {
		routesList[i] = *routesMap[routesList[i].Id]
	}

	return routesList, nil
}

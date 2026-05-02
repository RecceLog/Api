package services

import (
	"Api/internal/domain"
	"Api/internal/infrastructure/repositories"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoutesService interface {
	CreateWithNotes(ctx context.Context, r *domain.Route, ns *domain.NoteSet) error
	AddNoteSet(ctx context.Context, ns *domain.NoteSet, routeId uuid.UUID) error
	FindById(ctx context.Context, id uuid.UUID) (domain.Route, error)
	FindAll(ctx context.Context) ([]domain.Route, error)
	FindInRange(ctx context.Context, point domain.Coordinate, rangeM int) ([]domain.Route, error)
	FindNoteSets(ctx context.Context, routeId uuid.UUID) ([]domain.NoteSet, error)
	FindNoteSet(ctx context.Context, setId uuid.UUID) (domain.NoteSet, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteNoteSet(ctx context.Context, setId uuid.UUID) error
}

type routesService struct {
	pool       *pgxpool.Pool
	routesRepo repositories.RoutesRepository
	notesRepo  repositories.NotesRepository
}

func NewRoutesService(
	conn *pgxpool.Pool,
	routesRepo repositories.RoutesRepository,
	notesRepo repositories.NotesRepository,
) RoutesService {
	return &routesService{
		pool:       conn,
		routesRepo: routesRepo,
		notesRepo:  notesRepo,
	}
}

// CreateWithNotes inserts a route and its first note set atomically.
// If either insert fails, the whole operation is rolled back.
func (s *routesService) CreateWithNotes(ctx context.Context, r *domain.Route, ns *domain.NoteSet) error {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	r.Id, err = uuid.NewV7()
	if err != nil {
		return err
	}

	if err = s.routesRepo.Insert(ctx, tx, r); err != nil {
		return err
	}

	if err = s.notesRepo.InsertSet(ctx, tx, ns, r.Id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// AddNoteSet inserts an additional note set for an existing route.
func (s *routesService) AddNoteSet(ctx context.Context, ns *domain.NoteSet, routeId uuid.UUID) error {
	return s.notesRepo.InsertSet(ctx, s.pool, ns, routeId)
}

func (s *routesService) FindById(ctx context.Context, id uuid.UUID) (domain.Route, error) {
	return s.routesRepo.FindById(ctx, s.pool, id)
}

func (s *routesService) FindAll(ctx context.Context) ([]domain.Route, error) {
	return s.routesRepo.FindAll(ctx, s.pool)
}

func (s *routesService) FindInRange(ctx context.Context, point domain.Coordinate, rangeM int) ([]domain.Route, error) {
	return s.routesRepo.FindInRange(ctx, s.pool, point, rangeM)
}

func (s *routesService) FindNoteSets(ctx context.Context, routeId uuid.UUID) ([]domain.NoteSet, error) {
	return s.notesRepo.FindSetsByRouteId(ctx, s.pool, routeId)
}

func (s *routesService) FindNoteSet(ctx context.Context, setId uuid.UUID) (domain.NoteSet, error) {
	return s.notesRepo.FindSetById(ctx, s.pool, setId)
}

// Delete removes a route and all its associated note sets atomically.
func (s *routesService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.routesRepo.Delete(ctx, s.pool, id)
}

func (s *routesService) DeleteNoteSet(ctx context.Context, setId uuid.UUID) error {
	return s.notesRepo.DeleteSet(ctx, s.pool, setId)
}

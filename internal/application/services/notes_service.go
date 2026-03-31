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

type NotesService interface {
	InsertSet(ctx context.Context, noteSet *domain.NoteSet, routeId uuid.UUID) error
	FindSetsByRouteId(ctx context.Context, routeId uuid.UUID) ([]domain.NoteSet, error)
	FindSetById(ctx context.Context, setId uuid.UUID) (domain.NoteSet, error)
	DeleteSet(ctx context.Context, setId uuid.UUID) error
}

type notesService struct {
	connPool *pgxpool.Pool
}

func NewNotesService(connPool *pgxpool.Pool) NotesService {
	return &notesService{
		connPool: connPool,
	}
}

func (s *notesService) InsertSet(ctx context.Context, noteSet *domain.NoteSet, routeId uuid.UUID) error {

	// Generate note set id
	noteSetId, err := uuid.NewV7()
	if err != nil {
		slog.Error("Error creating UUID locally for note set", "error", err.Error())
		return err
	}
	noteSet.Id = noteSetId
	noteSet.RouteId = routeId

	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		slog.Error("Error creating transaction", "error", err.Error())
		return err
	}
	defer tx.Rollback(ctx)

	// Insert new note set
	_, err = tx.Exec(
		ctx,
		database.InsertNoteSet,
		noteSet.Id, noteSet.RouteId,
	)
	if err != nil {
		return err
	}

	// Insert every note for the set
	for i, note := range noteSet.Notes {
		// Generate uuid for note
		noteId, err := uuid.NewV7()
		if err != nil {
			slog.Error("Error creating UUID locally for note", "error", err.Error())
			return err
		}
		noteSet.Notes[i].Id = noteId

		// Prepare query to insert note data in db
		_, err = tx.Exec(
			ctx,
			database.InsertNote,
			noteId, noteSetId, note.Position.Lng,
			note.Position.Lat, i+1, note.Type, note.Severity,
			note.Direction, note.Description,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *notesService) FindSetsByRouteId(ctx context.Context, routeId uuid.UUID) ([]domain.NoteSet, error) {

	// Get notes rows from db
	rows, err := s.connPool.Query(
		ctx,
		database.GetNoteSetsByRouteId,
		routeId,
	)
	if err != nil {
		return nil, err
	}

	// Map rows to a list of notes
	notes, err := scanNotesRows(rows)

	return notes, err
}

func (s *notesService) FindSetById(ctx context.Context, setId uuid.UUID) (domain.NoteSet, error) {

	// Get notes rows from db
	rows, err := s.connPool.Query(
		ctx,
		database.GetSetById,
		setId,
	)
	if err != nil {
		return domain.NoteSet{}, err
	}

	// Map rows to a list of notes
	notes, err := scanNotesRows(rows)

	if len(notes) == 0 {
		return domain.NoteSet{}, pgx.ErrNoRows
	}

	return notes[0], err
}

func (s *notesService) DeleteSet(ctx context.Context, setId uuid.UUID) error {

	_, err := s.connPool.Exec(
		ctx,
		database.DeleteNoteSetById,
		setId,
	)
	return err
}

func scanNotesRows(rows pgx.Rows) ([]domain.NoteSet, error) {

	noteSetsMap := make(map[uuid.UUID]*domain.NoteSet)

	for rows.Next() {
		var noteSetId, routeId, noteId uuid.UUID
		var note domain.Note

		err := rows.Scan(
			&noteSetId, &routeId, &noteId,
			&note.Position.Lat, &note.Position.Lng,
			&note.Order, &note.Type, &note.Severity,
			&note.Direction, &note.Description,
		)
		if err != nil {
			return nil, err
		}

		// Set the note ID
		note.Id = noteId

		// Check if this NoteSet already exists in the map
		if noteSet, exists := noteSetsMap[noteSetId]; exists {
			// Append to existing NoteSet
			noteSet.Notes = append(noteSet.Notes, note)
		} else {
			// Create new NoteSet
			noteSetsMap[noteSetId] = &domain.NoteSet{
				Id:      noteSetId,
				RouteId: routeId,
				Notes:   []domain.Note{note},
			}
		}
	}

	// Convert map to slice
	var noteSets []domain.NoteSet
	for _, noteSet := range noteSetsMap {
		noteSets = append(noteSets, *noteSet)
	}

	return noteSets, nil
}

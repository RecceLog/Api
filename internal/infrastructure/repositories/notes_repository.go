package repositories

import (
	"Api/internal/domain"
	"Api/internal/infrastructure/database"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type NotesRepository interface {
	InsertSet(ctx context.Context, q database.Querier, noteSet *domain.NoteSet, routeId uuid.UUID) error
	FindSetsByRouteId(ctx context.Context, q database.Querier, routeId uuid.UUID) ([]domain.NoteSet, error)
	FindSetById(ctx context.Context, q database.Querier, setId uuid.UUID) (domain.NoteSet, error)
	DeleteSet(ctx context.Context, q database.Querier, setId uuid.UUID) error
}

type notesRepository struct{}

func NewNotesRepository() NotesRepository {
	return &notesRepository{}
}

func (r *notesRepository) InsertSet(ctx context.Context, q database.Querier, noteSet *domain.NoteSet, routeId uuid.UUID) error {

	// Generate note set id
	noteSetId, err := uuid.NewV7()
	if err != nil {
		slog.Error("Error creating UUID locally for note set", "error", err.Error())
		return err
	}
	noteSet.Id = noteSetId
	noteSet.RouteId = routeId

	// Insert new note set
	_, err = q.Exec(
		ctx,
		database.InsertNoteSet,
		noteSet.Id, noteSet.RouteId,
	)
	if err != nil {
		return err
	}
	replacer := strings.NewReplacer(
		"$1", noteSet.Id.String(),
		"$2", fmt.Sprintf("%s", routeId.String()),
	)
	fmt.Printf("%s;", replacer.Replace(database.InsertNoteSet))

	// Insert every note for the set
	for i, note := range noteSet.Notes {

		// Generate uuid for note
		noteId, err := uuid.NewV7()
		if err != nil {
			slog.Error("Error creating UUID locally for note", "error", err.Error())
			return err
		}
		noteSet.Notes[i].Id = noteId

		_, err = q.Exec(
			ctx,
			database.InsertNote,
			noteId, noteSetId, note.Position.Lng,
			note.Position.Lat, i+1, note.Type, note.Severity,
			note.Direction, note.Description,
		)
		if err != nil {
			return err
		}

		waypointReplacer := strings.NewReplacer(
			"$1", noteId.String(),
			"$2", fmt.Sprintf("%s", noteSet.Id.String()),
			"$3", fmt.Sprintf("%f", note.Position.Lng),
			"$4", fmt.Sprintf("%f", note.Position.Lat),
			"$5", fmt.Sprintf("%d", i+1),
			"$6", fmt.Sprintf("%s", note.Type),
			"$7", fmt.Sprintf("%d", note.Severity),
			"$8", fmt.Sprintf("%s", note.Direction),
			"$9", fmt.Sprintf("%s", note.Description),
		)
		fmt.Printf("%s;\n", waypointReplacer.Replace(database.InsertNote))
	}

	return nil
}

func (r *notesRepository) FindSetsByRouteId(ctx context.Context, q database.Querier, routeId uuid.UUID) ([]domain.NoteSet, error) {

	rows, err := q.Query(ctx, database.GetNoteSetsByRouteId, routeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNotesRows(rows)
}

func (r *notesRepository) FindSetById(ctx context.Context, q database.Querier, setId uuid.UUID) (domain.NoteSet, error) {

	rows, err := q.Query(ctx, database.GetSetById, setId)
	if err != nil {
		return domain.NoteSet{}, err
	}
	defer rows.Close()

	noteSets, err := scanNotesRows(rows)
	if err != nil {
		return domain.NoteSet{}, err
	}

	if len(noteSets) == 0 {
		return domain.NoteSet{}, pgx.ErrNoRows
	}

	return noteSets[0], nil
}

func (r *notesRepository) DeleteSet(ctx context.Context, q database.Querier, setId uuid.UUID) error {

	_, err := q.Exec(ctx, database.DeleteNoteSetById, setId)
	return err
}

func scanNotesRows(rows pgx.Rows) ([]domain.NoteSet, error) {

	noteSetsMap := make(map[uuid.UUID]*domain.NoteSet)
	var orderedIds []uuid.UUID

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

		note.Id = noteId

		if noteSet, exists := noteSetsMap[noteSetId]; exists {
			noteSet.Notes = append(noteSet.Notes, note)
		} else {
			noteSetsMap[noteSetId] = &domain.NoteSet{
				Id:      noteSetId,
				RouteId: routeId,
				Notes:   []domain.Note{note},
			}
			orderedIds = append(orderedIds, noteSetId)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Preserve insertion order (map iteration is random in Go)
	result := make([]domain.NoteSet, 0, len(orderedIds))
	for _, id := range orderedIds {
		result = append(result, *noteSetsMap[id])
	}

	return result, nil
}

package domain

import (
	"github.com/google/uuid"
)

type Note struct {
	Id          uuid.UUID  `json:"id,omitempty"`
	Position    Coordinate `json:"position" binding:"required"`
	Order       int        `json:"order"`
	Type        string     `json:"type" binding:"required"`
	Severity    int        `json:"severity"`
	Direction   string     `json:"direction"`
	Description string     `json:"description"`
}

type NoteSet struct {
	Id      uuid.UUID `json:"id"`
	RouteId uuid.UUID `json:"route_id"`
	Notes   []Note    `json:"notes" binding:"required,dive"`
}

package domain

import (
	"github.com/google/uuid"
)

type Route struct {
	Id        uuid.UUID        `json:"id"`
	Start     Coordinate       `json:"start" binding:"required"`
	Waypoints []Waypoint       `json:"waypoints" binding:"dive"`
	Finish    Coordinate       `json:"finish" binding:"required"`
}

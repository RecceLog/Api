package domain

type Waypoint struct {
	Position Coordinate `json:"position" binding:"required"`
	Order uint          `json:"order"`
}

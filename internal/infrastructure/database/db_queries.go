package database

const (

	InsertRoute = `
		INSERT INTO routes(id, start, finish)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), ST_SetSRID(ST_MakePoint($4, $5), 4326))`

	InsertRouteWaypoint = `
		INSERT INTO waypoints(route_id, position, "order")
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4)`

	SelectAllRoutes = `
		SELECT
			r.id,
			ST_Y(r.start), ST_X(r.start),
			ST_Y(r.finish), ST_X(r.finish),
			ST_Y(w.position), ST_X(w.position)
		FROM routes r
		LEFT JOIN waypoints w ON r.id = w.route_id
		ORDER BY r.id, w.order`

	SelectRouteById = `
		SELECT
			r.id,
			ST_Y(r.start), ST_X(r.start),
			ST_Y(r.finish), ST_X(r.finish),
			ST_Y(w.position), ST_X(w.position)
		FROM routes r
		LEFT JOIN waypoints w ON r.id = w.route_id
		WHERE r.id = $1
		ORDER BY r.id, w.order`

	SelectRoutesInRange = `
		SELECT
			r.id,
			ST_Y(r.start), ST_X(r.start),
			ST_Y(r.finish), ST_X(r.finish),
			ST_Y(w.position), ST_X(w.position)
		FROM routes r
		LEFT JOIN waypoints w ON r.id = w.route_id
		WHERE
			ST_DWithin(
				start::geography,
				ST_MakePoint($1, $2), $3)
		ORDER BY r.id, w.order`

	DeleteRouteById = `
		DELETE FROM routes
		WHERE id = $1`



	InsertNoteSet = `
		INSERT INTO note_sets(id, route_id)
		VALUES ($1, $2)`

	InsertNote = `
		INSERT INTO notes(id, set_id, position, "order", "type", severity, direction, "description")
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5, $6, $7, $8, $9)`

	GetSetById = `
		SELECT
			ns.id, ns.route_id, n.id,
			ST_Y(position), ST_X(position),
			"order", type, severity,
			direction, description
		FROM notes n
		LEFT JOIN note_sets ns ON ns.id = n.set_id
		WHERE ns.id = $1
		ORDER BY ns.route_id, ns.id, n.order`

	GetNoteSetsByRouteId = `
		SELECT
			ns.id, ns.route_id, n.id,
			ST_Y(position), ST_X(position),
			"order", type, severity,
			direction, description
		FROM notes n
		LEFT JOIN note_sets ns ON ns.id = n.set_id
		WHERE ns.route_id = $1
		ORDER BY ns.route_id, ns.id, n.order`

	GetNotesByRouteIdAndSet = `
		SELECT
			id, route_id, "set", ST_Y(position), ST_X(position), "order", type, severity, direction, description
		FROM notes
		WHERE route_id = $1 AND "set" = $2
		ORDER BY route_id, "set", "order"`

	DeleteNoteSetById = `
		DELETE FROM note_sets
		WHERE id = $1`
)

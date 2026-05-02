-- +goose Up
-- +goose StatementBegin
-- CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    start GEOMETRY NOT NULL,
    finish GEOMETRY NOT NULL
);

CREATE TABLE IF NOT EXISTS waypoints (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    position GEOMETRY NOT NULL,
    "order" INT NOT NULL
);

CREATE TABLE IF NOT EXISTS note_sets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE
);

CREATE TYPE note_type AS ENUM('INDICATION', 'WARNING');
CREATE TYPE direction_type AS ENUM('LEFT', 'RIGHT', 'STRAIGHT', 'CHICANE');
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    set_id UUID NOT NULL REFERENCES note_sets(id) ON DELETE CASCADE,
    position GEOMETRY NOT NULL,
    "order" INT NOT NULL,
    "type" NOTE_TYPE NOT NULL,
    severity INT,
    direction DIRECTION_TYPE,
    "description" VARCHAR(255)
);

WITH route AS (
INSERT INTO routes(start, finish)
VALUES (
    ST_SetSRID(ST_MakePoint(7.287241, 44.806422), 4326),
    ST_SetSRID(ST_MakePoint(7.200493, 44.791587), 4326)
    )
    RETURNING id
    ),
    note_set AS (
INSERT INTO note_sets(route_id)
SELECT id FROM route
    RETURNING id, route_id
    ),
    waypoints_insert AS (
INSERT INTO waypoints(route_id, position, "order")
SELECT route_id, ST_SetSRID(ST_MakePoint(7.272475, 44.806649), 4326), 1 FROM note_set
UNION ALL
SELECT route_id, ST_SetSRID(ST_MakePoint(7.255800, 44.804868), 4326), 2 FROM note_set
    RETURNING id
    )
INSERT INTO notes(set_id, position, "order", "type", severity, direction, "description")
SELECT note_set.id, ST_SetSRID(ST_MakePoint(7.282342, 44.806329), 4326), 1, 'INDICATION'::note_type, 4, 'RIGHT'::direction_type, 'test' FROM note_set
UNION ALL
SELECT note_set.id, ST_SetSRID(ST_MakePoint(7.281546, 44.806480), 4326), 2, 'INDICATION'::note_type, 5, 'LEFT'::direction_type, 'test' FROM note_set;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes;
DROP TYPE IF EXISTS direction_type;
DROP TYPE IF EXISTS note_type;
DROP TABLE IF EXISTS note_sets;
DROP TABLE IF EXISTS waypoints;
DROP TABLE IF EXISTS routes;
-- +goose StatementEnd

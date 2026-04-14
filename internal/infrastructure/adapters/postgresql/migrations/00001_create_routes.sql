-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION postgis;

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

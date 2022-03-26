-- migrate:up
CREATE TABLE boards
(
	id         BIGSERIAL PRIMARY KEY,
	name       VARCHAR(100) NOT NULL,
	game_id    BIGINT,
	width      INT          NOT NULL DEFAULT 0,
	height     INT          NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX boards_name_uidx ON boards (name);
CREATE INDEX boards_game_id_idx ON boards (game_id) WHERE game_id IS NOT NULL;

-- migrate:down
DROP TABLE boards;

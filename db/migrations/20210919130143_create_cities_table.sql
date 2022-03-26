-- migrate:up
CREATE TABLE cities
(
	id         BIGSERIAL PRIMARY KEY,
	board_id   BIGINT       not null,
	name       VARCHAR(100) NOT NULL,
	x          INT          not null default 0,
	y          INT          not null default 0,
	upgrade_offered SMALLINT,
	immediate_point SMALLINT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX cities_board_id_idx ON cities (board_id);

-- migrate:down
DROP TABLE cities;

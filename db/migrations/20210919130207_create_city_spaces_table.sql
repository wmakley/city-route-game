-- migrate:up
CREATE TABLE city_spaces
(
	id                 BIGSERIAL PRIMARY KEY,
	city_id            BIGINT      not null,
	"order"            SMALLINT    not null default 1,
	space_type         SMALLINT    not null default 1,
	required_privilege SMALLINT    not null default 1,
	created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX city_spaces_city_id_position_uidx ON city_spaces (city_id, "order");


-- migrate:down
DROP TABLE city_spaces;

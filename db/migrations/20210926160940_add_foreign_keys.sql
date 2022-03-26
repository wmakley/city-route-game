-- migrate:up
ALTER TABLE city_spaces ADD FOREIGN KEY (city_id) REFERENCES cities(id) ON UPDATE cascade;
ALTER TABLE cities ADD FOREIGN KEY (board_id) REFERENCES boards(id) ON UPDATE cascade;

-- migrate:down
ALTER TABLE cities DROP CONSTRAINT cities_board_id_fkey;
ALTER TABLE city_spaces DROP CONSTRAINT city_spaces_city_id_fkey;

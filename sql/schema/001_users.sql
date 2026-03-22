-- +goose Up
CREATE TABLE users(
	id integer primary key,
	created_at timestamp not null,
	updated_at timestamp not null,
	name text unique
);

-- +goose Down
DROP TABLE users;

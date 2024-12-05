-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

    CREATE TABLE IF NOT EXISTS "UserSQLModel" (
			id SERIAL PRIMARY KEY,
			firstname VARCHAR(100),
			surname VARCHAR(100),
			password1 VARCHAR(100),
			email VARCHAR(100) UNIQUE
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE UserSQLModel
-- +goose StatementEnd

migrate-up:
	goose postgres "host=localhost user=postgres port=5432 dbname=testhub sslmode=disable" up

migrate-down:
	goose postgres "host=localhost user=postgres port=5432 dbname=testhub sslmode=disable" down
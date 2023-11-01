migrate-up:
	goose -dir ./migrations postgres "host=localhost user=postgres port=5432 dbname=testhub sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=localhost user=postgres port=5432 dbname=testhub sslmode=disable" down

migrate-reset:
	goose -dir ./migrations postgres "host=localhost user=postgres port=5432 dbname=testhub sslmode=disable" reset
export BANKING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" && migrate -database ${BANKING} -path ./db/migrations down

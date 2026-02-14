# Generate Go code from schemas
[group("backend")]
[working-directory: "backend"]
sqlc *command="generate":
	sqlc {{ command }}

# Run the backend server
[group("backend")]
[working-directory: "backend"]
run-backend:
    go run cmd/main.go

[group("backend")]
[working-directory: "backend"]
migrate-create seq:
	migrate create -ext sql -dir sql/migrations -seq {{seq}}

[group("deploy")]
run-db-container:
    docker run -p5432:5432 kopify-postgres

[group("deploy")]
build-db-container:
    docker build -t kopify-postgres -f backend/DB.Dockerfile

[group("deploy")]
build-api-container:
    docker build -t kopify-backend -f backend/API.Dockerfile

[group("deploy")]
run-api-container:
    docker run -p8080:8080 kopify-backend

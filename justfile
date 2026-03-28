[group("frontend")]
[working-directory: "frontend"]
run-frontend:
    bun run dev

[group("frontend")]
[working-directory: "frontend"]
check-frontend:
    bun run type-check

# Generate Go code from schemas
[group("backend")]
[working-directory: "backend"]
sqlc *command="generate":
	sqlc {{ command }}

# Run the backend server
[group("backend")]
[working-directory: "backend"]
run-backend:
    air

# Run the backend server (no live reload)
[group("backend")]
[working-directory: "backend"]
run-backend-static:
    go run cmd/main.go

[group("backend")]
[working-directory: "backend"]
migrate-create seq:
	migrate create -ext sql -dir sql/migrations -seq {{seq}}

[group("deploy")]
run-db-container:
    docker run -p5432:5432 koopify-postgres

[group("deploy")]
build-db-container:
    docker build -t koopify-postgres -f backend/DB.Dockerfile

[group("deploy")]
build-api-container:
    docker build -t koopify-backend -f backend/API.Dockerfile

[group("deploy")]
run-api-container:
    docker run -p8080:8080 koopify-backend

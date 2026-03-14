# Koopify

A sassy little web app to transfer funds from your wallet to mine and transfer snake oil from my non-existant inventory to yours!

## Local Deployment

Docker compose is used for local deployments.

```
$ docker-compose --version
Docker Compose version 5.0.2
```

Links to locally hosted services:
- koopify-web: [https://localhost:8000/](https://localhost:8000/)
- koopify-api: [https://localhost:8080/](https://localhost:8080/)
- koopify-db: TODO

TLS is used, as the same container images are used for the real deployment on
[https://koopify.piguy.nl](https://koopify.piguy.nl). The docker-compose
manifest spins up an alpine container to generate locally signed certificates.
The official deployment uses letsencrypt.

## AI usage disclaimer

AI has been used to design the front-end and produce the style sheets for the
frontend. It has also been used to alter parts of the code, be it replicating
schema from the backend into TypeScript types in the front-end, or making
adjustments to the UI/layout of the front-end page.

## Environment variables

### Backend (API)

Required secrets/config:

- `JWT_SECRET`: Secret used to sign JWTs. Must be a non-default, non-empty value.
- `PGDB`: Postgres connection string. Default: `postgres://postgres:postgres@localhost:5432/?sslmode=disable`
- `HOST_ADDR`: Bind address for the API server. Default: `:8080`
- `TLS_ENABLED`: Enable TLS when set to a truthy value (any non-empty value other than `0`).
- `TLS_CERT`: Path to TLS certificate file (required when `TLS_ENABLED` is on).
- `TLS_KEY`: Path to TLS private key file (required when `TLS_ENABLED` is on).

### Frontend (Web)

Required config:

- `VITE_API_BASE_URL`: Base URL for the API (e.g. `https://localhost:8080`).

### License

All the code in this repository is licensed under Apache-2.0 license.

# Koopify

## Development

Basic development information:

Protobuf files are used to define the API schema, these files generate Go structs for the backend and TS interfaces for the frontend. This provides a type-safe bridge to communicate between the frontend and the backend. A HTTP+JSON API is used for the actual communications, as using protobuf is probably not allowed.

The app is fully containerised for easy deployment on kubernetes, openshift or any other platform which supports OCI containers. For local development, podman or docker can be used. For orchestration, docker-compose is used.

The frontend uses Vue.js v3 and is an SPA. The backend uses Go, and the database in use is a postgresql database. sqlc is used as the ORM of choice.

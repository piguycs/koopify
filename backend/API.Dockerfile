FROM docker.io/library/golang:1.25.6-alpine3.23 as build

WORKDIR /build

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o main ./cmd


FROM scratch
COPY --from=build /build/main /main
EXPOSE 8080
CMD ["/main"]

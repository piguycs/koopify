FROM docker.io/library/golang:1.25.6-alpine3.23

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/main.go

EXPOSE 8080

CMD ["/build/main"]

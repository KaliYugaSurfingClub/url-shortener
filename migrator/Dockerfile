FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./../bin/migrator cmd/migrator/main.go

FROM alpine AS runner

WORKDIR /usr/local/bin

COPY --from=builder /usr/local/bin/migrator ./
COPY migrations migrations

ENV POSTGRES_URL="POSTGRES_URL was not set"

CMD ["./migrator"]


FROM golang:1.20-alpine AS builder

RUN apk --no-cache add bash git make gcc gettext

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates postgresql-client && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o migrate .

COPY ./migrations ./migrations

ENV POSTGRES_URL="postgres url was not set"

CMD ["./migrate"]

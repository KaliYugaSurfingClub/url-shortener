FROM golang:1.23

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates postgresql-client && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o link-service ./..

ENV CONFIG_PATH="/app/config/config.local.yaml"

CMD ["./link-service"]
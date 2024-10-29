FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./../bin/link-service cmd/link-service/main.go

FROM alpine AS runner

WORKDIR /usr/local/bin

COPY --from=builder /usr/local/bin/link-service ./
COPY config config

CMD ["./link-service"]


FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./../bin/app cmd/app/main.go

FROM alpine AS runner

WORKDIR /usr/local/bin

COPY --from=builder /usr/local/bin/app ./
COPY config config

ENV POSTGRES_URL="POSTGRES_URL was not set"

EXPOSE 8081

CMD ["./app"]


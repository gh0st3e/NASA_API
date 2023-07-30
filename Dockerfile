FROM golang:alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine as runner

WORKDIR /app
COPY --from=builder /src/main /app/
COPY .env /app/.env
COPY docs /app/docs

ENTRYPOINT ./main

CMD ["/app"]
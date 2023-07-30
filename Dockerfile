FROM golang:alpine AS builder

WORKDIR /src
COPY . .

RUN go mod download
RUN go build -o main cmd/main.go

FROM alpine as runner

WORKDIR /app
COPY --from=builder /src/main /app/
COPY .env /app/.env

ENTRYPOINT ./main

CMD ["/app"]
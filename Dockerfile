FROM golang:alpine AS builder

ADD . /src
RUN cd /src && go build -o main cmd/main.go

FROM alpine as runner

WORKDIR /app
COPY --from=builder /src/main /app/
COPY .env /app/.env

ENTRYPOINT ./main

CMD ["/app"]
build:
	go build -o cmd/main cmd/main.go
run: build
	 ./cmd/main
test:
	go test ./...
coverage:
	go test -coverprofile=coverage.out ./...
docker-run:
	docker-compose up --build

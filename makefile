run:
	go run cmd/main.go

test:
	go test -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out
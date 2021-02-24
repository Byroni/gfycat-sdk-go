run:
	go run cmd/main

test:
	go test -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out
run:
	go run main.go
	
cover:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o coverage.html

test:
	go test -v ./...
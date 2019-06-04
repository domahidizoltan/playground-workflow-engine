build:
	go mod download
	go test -coverprofile cov.tmp ./internal/stockprice
	rm cov.tmp
	go build -o mockserver cmd/main.go

build-image:
	docker build -t stockprice-mockserver .
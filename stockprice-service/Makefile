build:
	go mod download
	go test -coverprofile cov.tmp ./internal/stockprice
	rm cov.tmp
	go build -o service cmd/main.go

build-image:
	docker build -t stockprice-service .
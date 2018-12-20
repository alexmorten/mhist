test:
	go test ./... -timeout 10s

run:
	go run main/main.go

image-build:
	docker build -t mhist .

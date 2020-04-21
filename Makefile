.PHONY: build
build:
	go build -o ./bin/previewer ./cmd/main.go
	chmod +x ./bin/previewer

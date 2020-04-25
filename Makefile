.PHONY: build
build:
	go build -o ./bin/previewer ./cmd/main.go
	chmod +x ./bin/previewer

run:
	docker-compose -f ./build/docker/docker-compose.yml up

clean:
	docker-compose -f ./build/docker/docker-compose.yml down

run_and_build:
	docker-compose -f ./build/docker/docker-compose.yml up --build

test:
	docker-compose -f ./build/docker/docker-compose-tests.yml up --abort-on-container-exit --exit-code-from tests && \
	docker-compose -f ./build/docker/docker-compose-tests.yml down
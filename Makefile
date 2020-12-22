temporal-local-up:
	docker-compose up -d

temporal-local-down:
	docker-compose down

test-unit:
	go test -v -short -race ./...
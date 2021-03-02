API_KEY=????

temporal-local-up:
	docker-compose up -d

temporal-local-down:
	docker-compose down

check: build lint test-unit

test-unit:
	go test -v -short -race ./...

build:
	go build ./...

lint:
	golangci-lint run

tools:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.34.1

run-square:
	go run cmd/microservice_square/main.go

run-volume:
	go run cmd/microservice_volume/main.go

run-workflow:
	go run cmd/microservice_workflow/main.go

run-datadog-agent:
	docker run -d --name dd-agent -p 8126:8126 -v /var/run/docker.sock:/var/run/docker.sock:ro -v /proc/:/host/proc/:ro -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro -e DD_API_KEY=$(API_KEY) -e DD_SITE="datadoghq.com" gcr.io/datadoghq/agent:7

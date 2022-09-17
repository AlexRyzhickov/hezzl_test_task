.PHONY: run-nats
run-nats:
	@docker run -p 4222:4222 -ti nats:latest

.PHONY: run-db
run-db:
	@docker run \
		-d \
		-v `pwd`/db:/docker-entrypoint-initdb.d/ \
		--rm \
		-p 5432:5432 \
		--name db \
		-e POSTGRES_DB=backend \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		postgres:12

.PHONY: run-service
run-service:
	@export PORT=8080 && export SUBJECT=foo && go run ./cmd/main.go

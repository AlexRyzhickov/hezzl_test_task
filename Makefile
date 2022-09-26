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

.PHONY: run-redis
run-redis:
	@docker run -p 6379:6379 --name some-redis -d redis

.PHONY: run-clickhouse
run-clickhouse:
	@docker run -p 9000:9000 -d --name some-clickhouse-server --ulimit nofile=262144:262144 clickhouse/clickhouse-server

.PHONY: run-campaigns-service
run-campaigns-service:
	@export PORT=8080 && export SUBJECT=foo && go run ./campaigns_sevice/cmd/main.go

.PHONY: run-logs-writer
run-logs-writer:
	@export SUBJECT=foo && export CLICKHOUSE_ADDR="127.0.0.1:9000" && go run ./writer_logs/main.go

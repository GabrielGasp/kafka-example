setup:
	@docker compose up -d
	@sleep 5
	@docker compose exec kafka kafka-topics --create --topic=users --bootstrap-server=localhost:29092 --partitions=3

run-consumer:
	@go run cmd/consumer/**.go

run-producer:
	@go run cmd/producer/**.go

teardown:
	@docker compose down --remove-orphans
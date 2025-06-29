.PHONY: help
help:
	@echo "WARN - create .env file in ./dev/docker directory"
	@echo "Available commands:"
	@echo "	make app - Run app in docker"
	@echo "	make nt  - Load testing via ghz"
	@echo "	make test - Run tests"


.PHONY: app
app:
	sudo docker-compose -f ./dev/docker/docker-compose-app.yml up -d --build


.PHONY: nt
nt:
	ghz --config ./dev/ghz_config.yaml localhost:50051


.PHONY: test
test:
	go test ./tests -v

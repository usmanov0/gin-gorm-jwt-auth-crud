-include .env
export

DOCKER_COMPOSE_FILE=docker-compose.yaml

arg = $(filter-out $@,$(MAKECMDGOALS))

.PHONY: start
start:
	@echo "Start Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} up -d ${DOCKER_SERVICES}
	sleep 2
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

.PHONY: stop
stop:
	@echo "Stop Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} stop ${DOCKER_SERVICES}
	sleep 2
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

.PHONY: stop
rm: stop
	@echo "Remove Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} rm -v -f ${DOCKER_SERVICES}

.PHONY: mod-download
mod-download:
	@echo "Go mod download"
	sleep 2
	docker-compose exec app go mod download

.PHONY: mod-tidy
mod-tidy:
	@echo "Go mod tidy"
	sleep 2
	docker-compose exec app go mod tidy

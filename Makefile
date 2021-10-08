include tron.mk

DOCKER_COMPOSE_RUN ?= docker-compose -f scripts/docker-compose.yaml --project-name=tailer

.PHONY: compose-run
compose-run:
	${DOCKER_COMPOSE_RUN} up osm-server

.PHONY: compose-download-regions
compose-download-regions:
	${DOCKER_COMPOSE_RUN} up download-regions

.PHONY: compose-download-stylesheets
compose-download-stylesheets:
	${DOCKER_COMPOSE_RUN} up download-stylesheets

.PHONY: compose-down
compose-down:
	${DOCKER_COMPOSE_RUN} down

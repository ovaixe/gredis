# ================================================================= #
# HELPER
# ================================================================= #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

	
# ================================================================= #
# DEVELOPMENT
# ================================================================= #

## run/server: run the 'cmp/server' application
.PHONY: run/server
run/server:
	go run ./cmd/server

## build: build the 'cmp/server' application
.PHONY: build
build:
	@echo 'Building cmd/server...'
	go build -o ./bin ./cmd/server

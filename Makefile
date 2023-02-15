.PHONY: Build proto Go

protoc := docker run \
		  -it --rm \
		  -u $$(id -u):$$(id -g) \
		  -v `pwd`:/defs \
		  namely/protoc-all

build-proto: ## build proto from source ./internal/proto/service.proto
	$(protoc)  -f ./internal/proto/service.proto -l go -o ./internal/proto/gen --go_opt=paths=source_relative
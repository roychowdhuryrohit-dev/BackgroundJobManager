.PHONY: all build_local run_local local build_docker clean_docker run_docker docker

all: local
# conf_env:
# 	. ./.env
build_local:
	go build -o bin/DemoService
run_local:
	. ./.env && ./bin/DemoService
local: build_local run_local
build_docker:
	docker build -t demo-service . && \
	docker images demo-service
run_docker:
	. ./.env && \
	docker run --env HOST --env GIN_MODE -p $$PORT:$$PORT demo-service 
clean_docker:
	docker container kill $$(docker container ls -aq); \
	docker rm $$(docker ps -a -q) \
	docker system prune --all --force --volumes; \
	docker image rm $$(docker image ls -a -q)
docker : build_docker run_docker
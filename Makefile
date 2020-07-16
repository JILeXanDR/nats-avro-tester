DOCKER_IMAGE_NAME=jilexandr/natsavrotester
GIT_TAG := $$(git describe --tag)

default: help

help:
	@echo "Available commands:"
	@echo "  server - starts dev server"
	@echo "  test - run all tests"
	@echo "  docker-server - start server inside docker container"
	@echo "  docker-build - build docker image"
	@echo "  docker-run - run docker image"

server:
	@echo "starting dev server..."
	go run .

docker-server:
	@echo "starting server inside docker container..."

docker-build:
	@echo "building docker image..."
	docker build -t ${DOCKER_IMAGE_NAME} -f Dockerfile .

docker-run:
	@echo "running docker image..."
	docker run -d --restart=always -p 9999:8080 ${DOCKER_IMAGE_NAME}

docker-compose:
	@echo "composing docker application..."
	PORT=9999 docker-compose up

docker-push:
	@echo "tagging an image ${DOCKER_IMAGE_NAME} using current git tag $(GIT_TAG)"
	@docker tag ${DOCKER_IMAGE_NAME} ${DOCKER_IMAGE_NAME}:$(GIT_TAG)
	@echo "pushing docker image ${DOCKER_IMAGE_NAME}:$(GIT_TAG) to Docker Hub..."
	docker push ${DOCKER_IMAGE_NAME}:${GIT_TAG}
	@echo "done!"

example:
	docker-compose up

test:
	@echo "running all tests..."
	@go test ./...
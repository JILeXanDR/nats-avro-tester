DOCKER_IMAGE_NAME=jilexandr/natsavrotester

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
	#docker build --no-cache -t ${DOCKER_IMAGE_NAME} -f Dockerfile .
	docker build -t ${DOCKER_IMAGE_NAME} -f Dockerfile .

docker-run:
	@echo "running docker image..."
	docker run -d --restart=always -p 9999:8080 ${DOCKER_IMAGE_NAME}

docker-compose:
	@echo "composing docker application..."
	PORT=9999 docker-compose up

docker-push:
	@echo "pushing docker image to Docker Hub..."
	docker push ${DOCKER_IMAGE_NAME}:latest

test:
	@echo "running all tests..."
	@go test ./...
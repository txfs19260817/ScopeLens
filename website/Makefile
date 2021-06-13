DOCKER_REPO = txfs19260817/scopelens-website

all: build push

build:
	docker build -t ${DOCKER_REPO} .

push:
	docker push ${DOCKER_REPO}
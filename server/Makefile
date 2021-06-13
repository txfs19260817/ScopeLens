DOCKER_REPO = txfs19260817/scopelens-server
# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

test:
	go test ./...

all: build push

build:
	docker build -t ${DOCKER_REPO} .

push:
	docker push ${DOCKER_REPO}

CONTAINER_IMAGE_TAG?=latest
CONTAINER_IMAGE_NAME?=quay.io/$(APP_ORGANIZATION)/$(APP_NAME):$(CONTAINER_IMAGE_TAG)

container-image: src/main/docker/Dockerfile.jvm pom.xml $(shell find src -type f -name '*')
	@echo Using Quarkus JVM build recipe
	mvn package
	docker build -f src/main/docker/Dockerfile.jvm -t $(CONTAINER_IMAGE_NAME) .
	echo $(CONTAINER_IMAGE_NAME) > container-image
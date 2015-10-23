# # Plain make targets if not requested inside a container
ifneq (,$(findstring test-integration,$(MAKECMDGOALS)))
	include Makefile.inc
	include mk/main.mk
else ifeq ($(USE_CONTAINER),)
	include Makefile.inc
	include mk/main.mk
else
# Otherwise, with docker, swallow all targets and forward into a container
DOCKER_IMAGE_NAME := "docker-machine-build"
DOCKER_CONTAINER_NAME := "docker-machine-build-container"
# get the dockerfile from docker/machine project so we stay in sync with the versions they use for go
DOCKER_FILE_URL := "https://raw.githubusercontent.com/docker/machine/master/Dockerfile"

build:
test: build
%:
		curl -s $DOCKER_FILE_URL > ./.dockerfile.machine
		docker build -f ./.dockerfile.machine -t $(DOCKER_IMAGE_NAME) .

		test -z '$(shell docker ps -a | grep $(DOCKER_CONTAINER_NAME))' || docker rm -f $(DOCKER_CONTAINER_NAME)

		docker run --name $(DOCKER_CONTAINER_NAME) \
		    -e DEBUG \
		    -e STATIC \
		    -e VERBOSE \
		    -e BUILDTAGS \
		    -e PARALLEL \
		    -e COVERAGE_DIR \
		    -e TARGET_OS \
		    -e TARGET_ARCH \
		    -e PREFIX \
		    $(DOCKER_IMAGE_NAME) \
		    make $@

		test ! -d bin || rm -Rf bin
		test -z "$(findstring build,$(patsubst cross,build,$@))" || docker cp $(DOCKER_CONTAINER_NAME):/go/src/github.com/docker/machine/bin bin

endif
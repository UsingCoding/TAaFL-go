DOCKER_BUILDER=docker buildx
DOCKER_BUILD=${DOCKER_BUILDER} build .

IMAGE_TAG?=master
export IMAGE?=vadimmakerov/compiler:${IMAGE_TAG}


all: build check

.PHONY: build
build: build-image

.PHONY: build-image
build-image:
	@${DOCKER_BUILD} \
	--target compiler-image \
	--output type=docker,name=${IMAGE} \
	--load # Load image from builder

# Build to docker tarball to chain-ability via pipes
.PHONY: build-docker-tarball
build-docker-tarball:
	@${DOCKER_BUILD} \
	--target compiler-image \
	--output type=docker,name=${IMAGE},dest=-

.PHONY: modules
modules:
	@${DOCKER_BUILD} \
 	--target go-mod-tidy \
	--output .

.PHONY: check
check:
	@${DOCKER_BUILD} \
 	--target lint

.PHONY: cache-clear
cache-clear: ## Clear the builder cache
	@${DOCKER_BUILDER} prune --force

PROJECTNAME ?= ksbagent
VERSION ?= 0.0.0

RANDOM := $(shell bash -c 'echo $$RANDOM')

.PHONY: build
build:
	docker build \
		--network=host \
		--build-arg VERSION=${VERSION} \
		-f build/Dockerfile \
		-t ${PROJECTNAME}builder .
	docker create --name ${PROJECTNAME}extract${RANDOM} ${PROJECTNAME}builder
	docker cp ${PROJECTNAME}extract${RANDOM}:/builded/. .builded
	docker rm -f ${PROJECTNAME}extract${RANDOM}

prune:
	docker container prune

remove:
	docker rmi $$(sudo docker images -qa -f 'dangling=true')

REPONAME = client
DOCKERIMAGENAME = benchflow/$(REPONAME)
VERSION = dev

.PHONY: all build_container_local 

all: build_container_local

install:
	wget -qO- https://github.com/benchflow/client/releases/download/v-dev/benchflow/getBenchFlow.sh | sh

build_container_local:
	docker build --no-cache -t $(DOCKERIMAGENAME):$(VERSION) -f Dockerfile.test .

test_container_local:
	docker run -ti --rm \
	--name $(REPONAME) $(DOCKERIMAGENAME):$(VERSION) python /app/benchflow.py
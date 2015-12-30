REPONAME = client
DOCKERIMAGENAME = benchflow/$(REPONAME)
VERSION = dev
GOPATH_SAVE_RESTORE:=`pwd`"/Godeps/_workspace"
GOPATH_BUILD:=`pwd`

.PHONY: all build_release 

all: build_release

save_dependencies:
	cd src/cloud/benchflow/$(REPONAME)/ && \
	GOPATH=$(GOPATH_SAVE_RESTORE) godep save ./... && \
	rm -rf ../../../../Godeps/*.* && \
	rm -rf ../../../../Godeps && \
	mv Godeps/ ../../../.. && \
	cd ../../../..

restore_dependencies: 
	GOPATH=$(GOPATH_SAVE_RESTORE) godep restore ./...

clean:
	go clean -i ./...
	rm -rf Godeps/_workspace/pkg

build:
	GOPATH=$(GOPATH_BUILD) godep go build -v ./...

build_release:
	GOPATH=$(GOPATH_BUILD) GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -ldflags '-s' -v ./...

install:
	godep go install -v ./...
	mv bin/$(REPONAME) bin/$(REPONAME)

test:
	godep go test ./...
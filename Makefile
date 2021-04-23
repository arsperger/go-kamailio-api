PKGNAME = go-kamailio-api
GITUSER = voip-services

GITREPO = gitlab.com/voip-services/go-kamailio-api

# CI_REGISTRY ?= registry.gitlab.com
# $(CI_JOB_TOCKEN) - deploy tocken

VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-list -1 HEAD)

LDFLAGS := -ldflags "-X $(GITREPO)/internal/utils.commit=$(COMMIT) \
	-X $(GITREPO)/internal/utils.version=$(VERSION)"

GCFLAGS := -gcflags='all=-N -l'

DOCKER_IMAGE = $(CI_REGISTRY)/$(GITUSER)/$(PKGNAME):$(VERSION)


clean:
		rm $(PKGNAME)

cleanall:
		rm go.mod
		rm go.sum
		rm $(PKGNAME)

lint:
		golint ./...

test:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -v -count 1 ./api/models -run Validate

testmain:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -v -count 1 -run Main

build: test
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(PKGNAME)

debug:
	   	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GCFLAGS) -o $(PKGNAME)

dockerbuild:
		docker build --build-arg GITREPO=$(GITREPO) --build-arg PROJECT=$(PKGNAME) --build-arg COMMIT=$(COMMIT) \
		--build-arg  VERSION=$(VERSION) -t $(DOCKER_IMAGE) -f Dockerfile .

dockerlogin:
		echo -n $(CI_JOB_TOKEN) | docker login -u gitlab-ci-token --password-stdin $(CI_REGISTRY)

release:
		docker push $(DOCKER_IMAGE)

.DEFAULT_GOAL := build

.PHONY: init build debug clean cleanall test testmain dockerbuild release

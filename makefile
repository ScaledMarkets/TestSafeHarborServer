# Makefile for building the tests for Safe Harbor Server.
# 


PRODUCTNAME=Safe Harbor Server
ORG=Scaled Markets
PACKAGENAME=testsafeharbor
EXECNAME=$(PACKAGENAME)
registryUser=testuser
registryPassword=testpassword
TestImageName=atomic
TestImageTag=atomicapp

.DELETE_ON_ERROR:
.ONESHELL:
.SUFFIXES:
.DEFAULT_GOAL: all

SHELL = /bin/sh

CURDIR=$(shell pwd)

.PHONY: all compile clean info
.DEFAULT: all

src_dir = $(CURDIR)/src

build_dir = $(CURDIR)/bin

all: compile

$(build_dir):
	mkdir $(build_dir)

$(build_dir)/$(EXECNAME): $(build_dir) $(src_dir)

# 'make compile' builds the executable, which is placed in <build_dir>.
compile: $(build_dir)/$(PACKAGENAME)

$(build_dir)/$(PACKAGENAME): src/..
	@GOPATH=$(CURDIR) go install $(PACKAGENAME)

setup:
	mkdir -p auth
	mkdir -p registrydata
	sudo docker run --entrypoint htpasswd docker.io/registry:2 -Bbn $(registryUser) $(registryPassword) > auth/htpasswd
	sudo docker run --net=host -d -p 5000:5000 --name registry \
		-v registryauth:/auth \
		-v registrydata:/var/lib/registry \
		-e "REGISTRY_AUTH=htpasswd" \
		-e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
		-e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
		docker.io/registry:2
	sudo docker pull docker.io/projectatomic/$(TestImage) && docker tag docker.io/projectatomic/$(TestImage) localhost:5000/$(TestImage)
	sudo docker push localhost:5000/$(TestImage)

run:
	bin/testsafeharbor

clean:
	rm -r -f $(build_dir)/$(PACKAGENAME)

info:
	@echo "Makefile for $(PRODUCTNAME)"


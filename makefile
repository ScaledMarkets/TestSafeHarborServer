# Makefile for building the tests for Safe Harbor Server.
# 

PRODUCTNAME=Safe Harbor Server
ORG=Scaled Markets
PACKAGENAME=testsafeharbor
EXECNAME=$(PACKAGENAME)

# These are needed by the registry tests:
registryUser=testuser
registryPassword=testpassword
TestImageName=atomicapp
TestImageTag=latest

# Needed by the SafeHarbor tests:
SAFEHARBOR_PORT=6000

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

# This target can only be run on a Linux system that has docker-engine installed.
prepregistry:
	# Create directories needed by the docker registry.
	mkdir -p registryauth
	mkdir -p registrydata
	# Create htpassword file containing a user and password.
	sudo docker run --entrypoint htpasswd docker.io/registry:2 \
		-Bbn $(registryUser) $(registryPassword) > registryauth/htpasswd

# This target can only be run on a Linux system that has docker-engine installed.
startregistry:
	# Start a docker registry instance.
	sudo docker rm -f registry
	sudo docker run -d -p 5000:5000 --name registry \
		-v `pwd`/registryauth:/auth \
		-v `pwd`/registrydata:/var/lib/registry \
		-e "REGISTRY_AUTH=htpasswd" \
		-e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
		-e "REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd" \
		docker.io/registry:2
		
stopregistry:
	sudo docker stop registry

# This target can only be run on a Linux system that has docker-engine installed.
getatomicapp:
	# Pull atomicapp to our docker client.
	sudo docker pull docker.io/projectatomic/$(TestImageName)
	sudo docker tag docker.io/projectatomic/$(TestImageName) localhost:5000/$(TestImageName)
	# Push atomic to our registry.
	sudo docker push localhost:5000/$(TestImageName)

runall:
	bin/testsafeharbor \
		h=$(SAFEHARBOR_HOST) \
		p=$(SAFEHARBOR_PORT) \
		-redispswd=ahdal8934k383898&*kdu&^ \
		-tests="Registry,json,goredis,redis,CreateRealmsAndUsers,CreateResources,CreateGroups,GetMy,AccessControl,UpdateAndReplace,Delete,DockerFunctions"

run:
	bin/testsafeharbor \
		-tests="Registry"

clean:
	rm -r -f $(build_dir)/$(PACKAGENAME)

info:
	@echo "Makefile for $(PRODUCTNAME). E.g.: make SAFEHARBOR_HOST=127.0.0.1 runall"


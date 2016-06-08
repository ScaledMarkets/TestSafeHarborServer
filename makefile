# Makefile for building the tests for Safe Harbor Server.

SHHOST=52.24.161.169
SHPORT=6000

PRODUCTNAME=Safe Harbor Server
ORG=Scaled Markets
PACKAGENAME=testsafeharbor
EXECNAME=$(PACKAGENAME)

# These are needed by the registry tests:
RegistryHost=localhost
RegistryPort=5000
registryUser=testuser
registryPassword=testpassword
TestImageRepoName=booploink
TestImageTag=latest
ImageToUploadPath=booploink
BooPloinkImageDigest=d2cf21381ce5a17243ec11062b5df136a9d5eac40c7bcdb3f65f42b32342c802
ImageToUploadDigest=$(BooPloinkImageDigest)

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

# Shortcut task for stopping, cleaning up, and restarting. Run this task after
# starting docker. After this task, testing tasks are ready to run.
testprep: stopregistry cleanregistry prepregistry startregistry

$(build_dir):
	mkdir $(build_dir)

$(build_dir)/$(EXECNAME): $(build_dir) $(src_dir)

# 'make compile' builds the executable, which is placed in <build_dir>.
compile: $(build_dir)/$(PACKAGENAME)

$(build_dir)/$(PACKAGENAME): $(build_dir)
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
	sudo docker run -d -p $(RegistryPort):$(RegistryPort) --name registry --net=host \
		-v `pwd`/registryauth:/auth \
		-v `pwd`/registrydata:/var/lib/registry \
		-e "REGISTRY_AUTH=htpasswd" \
		-e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
		-e "REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd" \
		docker.io/registry:2
		
stopregistry:
	sudo docker stop registry

cleanregistry:
	rm -r registryauth
	rm -r registrydata

# This target can only be run on a Linux system that has docker-engine installed.
getatomicapp:
	# Pull atomicapp to our docker client.
	sudo docker pull docker.io/projectatomic/atomicapp
	sudo docker tag docker.io/projectatomic/atomicapp $(RegistryHost):$(RegistryPort)/atomicapp
	# Push atomic to our registry.
	sudo docker login -u=$(registryUser) -p=$(registryPassword) -e="" $(RegistryHost):$(RegistryPort)
	sudo docker push $(RegistryHost):$(RegistryPort)/atomicapp

# Run all SafeHarborServer tests.
runall:
	bin/testsafeharbor -stop \
		-h=$(SHHOST) -p=$(SHPORT) \
		-tests="CreateRealmsAndUsers,CreateResources,CreateGroups,GetMy,AccessControl,UpdateAndReplace,Delete"
#		-tests="CreateRealmsAndUsers,CreateResources,CreateGroups,GetMy,AccessControl,UpdateAndReplace,Delete,DockerFunctions"

# Run redis tests.
runredis:
	bin/testsafeharbor \
		-redispswd="ahdal8934k383898&*kdu&^"
		-tests="redis"

# Unit test the DockerRegistry module.
regtests:
	export RegistryHost=$(RegistryHost)
	export RegistryPort=$(RegistryPort)
	export registryUser=$(registryUser)
	export registryPassword=$(registryPassword)
	export TestImageRepoName=$(TestImageRepoName)
	export TestImageTag=$(TestImageTag)
	export ImageToUploadPath=$(ImageToUploadPath)
	export ImageToUploadDigest=$(ImageToUploadDigest)
	bin/testsafeharbor -stop \
		-tests="Registry"

# Unit test the DockerEngine module.
engtests:
	export RegistryHost=$(RegistryHost)
	export RegistryPort=$(RegistryPort)
	export registryUser=testuser
	export registryPassword=testpassword
	bin/testsafeharbor -stop \
		-tests="Engine"

# Unit test the DockerServices module.
svctests:
	bin/testsafeharbor -stop \
		-h=$(SHHOST) -p=$(SHPORT) \
		-tests="DockSvcs"

# Run SafeHarborServer docker related tests.
dockertests:
	bin/testsafeharbor -stop \
		-h=$(SHHOST) -p=$(SHPORT) \
		-tests="DockerFunctions"

scanconfigs:
	bin/testsafeharbor -stop \
		-h=$(SHHOST) -p=$(SHPORT) \
		-tests="ScanConfigs"

# Run a SafeHarborServer "smoke" test suite.
basic:
	bin/testsafeharbor -stop \
		-h=$(SHHOST) -p=$(SHPORT) \
		-tests="CreateRealmsAndUsers"

# Run the SafeHarborServer "Delete" test suite.
delete:
	bin/testsafeharbor -stop \
		-h=$(SHHOST) -p=$(SHPORT) \
		-tests="Delete"

# List the images in the registry.
listimages:
	curl http://$(registryUser):$(registryPassword)@$(RegistryHost):$(RegistryPort)/v2/_catalog

# List the registry tags for the test image ("repo").
checkimage:
	curl http://$(registryUser):$(registryPassword)@$(RegistryHost):$(RegistryPort)/v2/$(TestImageRepoName)/tags/list

# Remove compilation artifacts.
clean:
	rm -r -f $(build_dir)/$(PACKAGENAME)

info:
	@echo "Makefile for $(PRODUCTNAME). E.g.: make SAFEHARBOR_HOST=127.0.0.1 runall"


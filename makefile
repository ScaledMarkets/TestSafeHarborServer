# Makefile for building the tests for Safe Harbor Server.
# This does not run any tests: it merely complies the code.


PRODUCTNAME=Safe Harbor Server
ORG=Scaled Markets
PACKAGENAME=testsafeharbor
EXECNAME=$(PACKAGENAME)

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

clean:
	rm -r -f $(build_dir)/$(PACKAGENAME)

info:
	@echo "Makefile for $(PRODUCTNAME)"


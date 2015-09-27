# Makefile for building the tests for Safe Harbor Server.
# This does not run any tests: it merely complies the code.


PRODUCTNAME=Safe Harbor Server
ORG=Scaled Markets
EXECNAME=TestSafeHarborServer

.DELETE_ON_ERROR:
.ONESHELL:
.SUFFIXES:
.DEFAULT_GOAL: all

SHELL = /bin/sh

CURDIR=$(shell pwd)

#GO_LDFLAGS=-ldflags "-X `go list ./version`.Version $(VERSION)"

.PHONY: all compile test clean info
.DEFAULT: all

src_dir = $(CURDIR)/src

build_dir = $(CURDIR)/../bin

GOPATH = $(CURDIR)

all: compile test

compile:
	@echo GOPATH=$(GOPATH)
	GOPATH=$(CURDIR) go build -o $(build_dir)/testmain main

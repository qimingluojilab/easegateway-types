.PHONY: default build fmt vendor_clean vendor_get vendor_update

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))

GOPATH := ${MKFILE_DIR}_vendor:${MKFILE_DIR}
export GOPATH

GATEWAY_TYPES_ALL_SRC_FILES = $(shell find ${MKFILE_DIR} -type f -name "*.go")

default: ${GATEWAY_TYPES_ALL_SRC_FILES}
	@echo "-------------- building gateway types --------------"
	echo ${GOPATH}
	cd ${MKFILE_DIR} && go build  -gcflags "-N -l" -v ./...

build: default

fmt:
	cd ${MKFILE_DIR} && go fmt ./...

vendor_clean:
	rm -dRf ${MKFILE_DIR}_vendor/src

vendor_get:
	GOPATH=${MKFILE_DIR}_vendor go get -d -u -v \
		github.com/hexdecteam/easegateway-types/...

vendor_update: vendor_get
	cd ${MKFILE_DIR} && rm -rf `find ./_vendor/src -type d -name .git` \
	&& rm -rf `find ./_vendor/src -type d -name .hg` \
	&& rm -rf `find ./_vendor/src -type d -name .bzr` \
	&& rm -rf `find ./_vendor/src -type d -name .svn`

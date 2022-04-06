# The development version of clang is distributed as the 'clang' binary,
# while stable/released versions have a version number attached.
# Pin the default clang to a stable version.
CLANG ?= clang-14
STRIP ?= llvm-strip-14
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)

# Obtain an absolute path to the directory of the Makefile.
# Assume the Makefile is in the root of the repository.
REPODIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Prefer podman if installed, otherwise use docker.
# Note: Setting the var at runtime will always override.
CONTAINER_ENGINE ?= docker

IMAGE := quay.io/cilium/ebpf-builder
VERSION := 1648566014

#.PHONY: all clean container-all container-shell generate

.DEFAULT_TARGET = container-all

# Build all ELF binaries using a containerized LLVM toolchain.
container-all:
	${CONTAINER_ENGINE} run --rm \
		-v "${REPODIR}":/ebpf -w /ebpf --env MAKEFLAGS \
		--env CFLAGS="-fdebug-prefix-map=/ebpf=." \
		--env HOME="/tmp" \
		"${IMAGE}:${VERSION}" \
		$(MAKE) all

# (debug) Drop the user into a shell inside the container as root.
container-shell:
	${CONTAINER_ENGINE} run --rm -ti \
		-v "${REPODIR}":/ebpf -w /ebpf \
		"${IMAGE}:${VERSION}"

clean:
	sudo rm ./bpf_bpfel.o ./bpf_bpfel.go ./main
	sudo rm -rf /sys/fs/bpf/ttcp_filter_table

all: generate

# $BPF_CLANG is used in go:generate invocations.
generate: export BPF_CLANG := $(CLANG)
generate: export BPF_CFLAGS := $(CFLAGS)
generate:
	go generate .

local:
	go build -o main *.go
	sudo ./main

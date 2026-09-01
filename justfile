set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

image := "dungar-dungar:latest"
platform := "linux/arm64"
archive := "dungar-arm64.tar"

default:
  @just --list

# Download and verify Go module dependencies.
dep:
  go mod download
  go mod verify

# Build the local dungar binary.
build: dep
  go build -mod=mod -o bin/dungar -compiler gc -ldflags="-s -w" cmd/dungar/*.go

# Build the ARM64 Docker image into the local Docker image store.
docker-image-arm64:
  docker buildx build --platform {{platform}} -t {{image}} --load .

# Rebuild dungar-arm64.tar from the ARM64 Docker image.
tarball: docker-image-arm64
  docker save -o {{archive}} {{image}}

# Build the ARM64 Docker tarball directly without loading the image locally.
tarball-direct:
  docker buildx build --platform {{platform}} -t {{image}} --output type=docker,dest={{archive}} .

# Start the local Docker Compose stack.
compose-up:
  docker compose up --build

# Stop the local Docker Compose stack.
compose-down:
  docker compose down

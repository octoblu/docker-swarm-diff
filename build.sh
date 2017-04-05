#!/bin/bash

APP_NAME=docker-swarm-diff
BUILD_DIR=$PWD/dist
IMAGE_NAME=local/$APP_NAME

build_on_local() {
  local goos="$1"
  local goarch="$2"
  local extension=""

  if [ "$goos" == "windows" ]; then
    extension=".exe"
  fi

  env CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" go build -a -tags netgo -installsuffix cgo -ldflags '-w' -o "${BUILD_DIR}/${APP_NAME}-${goos}-${goarch}${extension}" .
}

init() {
  rm -rf "${BUILD_DIR:?}/" \
  && mkdir -p "${BUILD_DIR:?}/"
}

package() {
  docker build --tag $IMAGE_NAME:latest entrypoint
}

fatal() {
  local message=$1
  echo $message
  exit 1
}

cross_compile_build(){
  for goos in darwin linux; do
    for goarch in 386 amd64; do
      build_on_local "$goos" "$goarch" > /dev/null
    done
  done
}

main() {
  local goos="$1"

  if [ -n "$goos" ]; then
    build_on_local "$goos" "amd64" > /dev/null
    exit $?
  fi

  cross_compile_build
}
main "$@"

#!/bin/bash

os=linux
arch=amd64
module=""
version=""
container=0
cwd=$(pwd)

for arg in "$@"
do
  case $arg in
    --module)
      shift
      module=$1
      shift
      ;;
    --version)
      shift
      version=$1
      shift
      ;;
    --image)
      container=1
      shift
      ;;
  esac
done

if [[ -z "$module" ]] || [[ ! -d "$module" ]] || [[ ! -f "$module/go.mod" ]] ; then
  echo "$module is not a valid module."
  exit 1
fi

bin=$module
build_root=build/
bin_path=$build_root/$bin

if [ -z $bin ]; then
  echo "A binary name must be specified through the module flag."
fi

if [ -z $version ]; then
  echo "A version must be specified."
  exit 1
fi

if [ -z $build_root ]; then
  exit 1
fi

cd $module

if [ -d $build_root ]; then
  rm -rf $build_root/*
fi

mkdir -p $build_root

go test ./...
if [ $? -ne 0 ]; then
  exit 1
fi

CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build \
  -o $bin_path \
  -ldflags="-s -w" \
  -trimpath main.go

if [ $container -eq 1 ] && [ "$os" == "linux" ]; then
  docker build -t $bin:$version --platform $os/$arch --build-arg BIN=$bin .
fi

cd $cwd

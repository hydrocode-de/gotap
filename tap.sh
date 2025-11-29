#!/bin/bash
set -e
go build -o tap .

if [ "$#" -gt 0 ] && [ -d "data/$1" ]; then
  cd "data/$1/src"
  shift
  ../../../tap "$@"
else
  cd data/valid/src
  ../../../tap "$@"
fi

cd ../../../
rm -f tap
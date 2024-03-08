#!/bin/bash

make build
result=$?

if [ $result -ne 0 ]; then
  echo "Make Build Failed."
  exit 1
fi

./bin/dungar learn1-file infra/alice-in-wonderland.txt
result=$?

if [ $result -ne 0 ]; then
  echo "Learn-File Failed!"
  exit 2
fi

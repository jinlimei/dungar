#!/bin/bash

#apt-get update && \
#  apt-get install -y netcat
#
#while ! nc -z postgres 5432;
#do
#  echo sleeping;
#  sleep 1;
#done;
#echo Connected!;

# Maybe not needed anymore! Yay!
#./tools/init_database.sh
#result=$?
#
#if [[ "$result" -ne 0 ]]; then
#  echo "INIT Database failed!"
#  exit 1
#fi
#
#./tools/seed_database.sh
#result=$?
#
#if [[ "$result" -ne 0 ]]; then
#  echo "SEED Database failed!"
#  exit 1
#fi

echo "LINTING"
make lint
result=$?

if [[ "$result" -ne 0 ]]; then
  echo "LINTING failed!"
  exit 1
fi

echo "MAKE TEST"
make test
result=$?

if [[ "$result" -ne 0 ]]; then
  echo "TEST failed!"
  exit 1
fi

echo "MAKE RACE"
make race
result=$?

if [[ "$result" -ne 0 ]]; then
  echo "RACE failed!"
  exit 1
fi

echo "MAKE MSAN"
make msan
result=$?

if [[ "$result" -ne 0 ]]; then
  echo "MSAN failed!"
  exit 1
fi

echo "MAKE BUILD"
make build
result=$?

if [[ "$result" -ne 0 ]]; then
  echo "BUILD failed!"
  exit 1
fi

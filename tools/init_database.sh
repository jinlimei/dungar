#!/bin/bash

PGPASSWORD=dungar_ci_test \
  psql -h postgres -U dungar_ci_test -d dungar_ci_test < ./infra/db/structure.sql

result=$?

if [[ $result -ne 0 ]]; then
  echo "Init Database Failed."
  exit 1
fi

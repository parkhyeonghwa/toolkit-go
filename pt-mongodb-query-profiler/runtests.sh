#!/bin/bash
export BASEDIR=$(git rev-parse --show-toplevel)/pt-mongodb-query-profiler
export CHECK_SESSIONS=0
cd $BASEDIR

go test -v -coverprofile=${BASEDIR}/coverage.out


if [ -f coverage.out ]
then
    go tool cover -func=coverage.out
    rm coverage.out
fi

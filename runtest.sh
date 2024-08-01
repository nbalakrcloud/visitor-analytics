#!/bin/bash

default="coverage.out"
coverage_path=${1:-$default}

if [ "$1" == "-coverage" ]; then
	mkdir -p coverage
	coverage_path="coverage/coverage.out"
fi

go test -v github.com/nbalakr/visitor-analytics/... -coverprofile=$coverage_path
exitcode=$?

echo Exit Code: $exitcode && exit $exitcode
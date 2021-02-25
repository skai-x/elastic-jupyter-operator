#!/usr/bin/env bash

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

cd ${SCRIPT_ROOT}
addlicense -f ./hack/license.txt *.go pkg/**/*.go
cd - >/dev/null

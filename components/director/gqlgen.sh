#!/usr/bin/env bash

echo "Generating code from GraphQL schema..."
COMPONENT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"

cd "$(dirname "$0")"

cd ${COMPONENT_DIR}/pkg/graphql/externalschema
GO111MODULE=on go run ../../../hack/gqlgen.go


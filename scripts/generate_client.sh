#!/bin/bash

rm -rf internal/spacetraders

openapi-generator-cli generate -i 'https://stoplight.io/api/v1/projects/spacetraders/spacetraders/nodes/reference/SpaceTraders.json\?fromExportButton\=true\&snapshotType\=http_service' \
    -o internal/api \
    -g go \
    --additional-properties=packageName="api" \
    --additional-properties=withGoMod=false

files=("go.mod" "go.sum" "git_push.sh" ".openapi-generator" "api" "test" ".travis.yml")

for f in ${files[@]}; do
    echo "internal/api/$f"
    rm -rf "internal/api/$f"
done

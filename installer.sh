#!/bin/bash

os=$(uname | tr '[:upper:]' '[:lower:]')
arch=""
owner="torana-us"
repo="tfdir"

if [ $(uname -m) = "x86_64" ]; then
    arch="amd64"
else
    arch="arm64"
fi

release_id=$(curl -H "Authorization: Bearer $GITHUB_TOKEN" \
    -sL https://api.github.com/repos/$owner/$repo/releases/latest \
    | jq '.assets[] | select(.name | contains("'$os'_'$arch'")) | .id')

curl -sLJ \
    -H 'Accept: application/octet-stream' \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    "https://api.github.com/repos/$owner/$repo/releases/assets/$release_id" -o tfdir.tar.gz

tar -xvf tfdir.tar.gz
rm tfdir.tar.gz

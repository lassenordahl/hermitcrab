#!/bin/bash

APP_DIR="$1"
PROJECT_NAME="$2"
VERSION="$3"

mkdir -p ./artifacts/ui
mkdir -p ./artifacts/outputs
cd "$APP_DIR"
pnpm run build
cp -r "dist/" ../artifacts/ui

# Tar and gzip the artifacts/ui directory
cd ../artifacts/ui
tar -zcf "../../artifacts/outputs/$PROJECT_NAME-ui-$VERSION.tar.gz" .
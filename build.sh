#!/usr/bin/env bash

#
# Build the whole project, before creating a proper Dockerfile
#

set -e

REQUIRED_FOLDER="wasaphoto"

if [ `basename $(pwd)` != "$REQUIRED_FOLDER" ]; then
	echo "Current directory is not $REQUIRED_FOLDER. Run this script in the appropriate folder."
	exit 1
fi

(
	cd webui/
	rm -rf node_modules/
	npm install
	npm run build-embed
)

go build -tags webui ./cmd/webapi

(
	cd webui/
	npm run clean
	git checkout HEAD -- node_modules/
)

echo "--------------------------------------"
echo "---- ./webapi executable created! ----"
echo "--------------------------------------"

#!/bin/sh
#
# Usage: deploy.sh TOKEN GITHUB_REPO GITRULES_RELEASE

gitrules -v github deploy --token=$1 --project=$2 --release=$3

#!/bin/sh
#
# Usage: erase.sh TOKEN GITHUB_REPO

set -e -x

gitrules -v github remove --token=$1 --repo=$2-gov.public
gitrules -v github remove --token=$1 --repo=$2-gov.private

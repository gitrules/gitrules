#!/bin/sh
#
# Usage: mkconfig.sh GITHUB_REPO ACCESS_TOKEN

GITHUB_REPO=$1
ACCESS_TOKEN=$2

CACHE_DIR=""
PUBLIC_REPO_URL="https://github.com/${GITHUB_REPO}-gitrules-public.git"
PRIVATE_REPO_URL="https://github.com/${GITHUB_REPO}-gitrules-private.git"

CONFIG_JSON=$(
     jq -n \
          --arg cache_dir "$CACHE_DIR" \
          --arg gov_pub_repo "$PUBLIC_REPO_URL" \
          --arg gov_priv_repo "$PRIVATE_REPO_URL" \
          --arg gov_auth_token "$ACCESS_TOKEN" \
          '{
               "cache_dir": $cache_dir,
               "auth" : {
                    ($gov_pub_repo): { "access_token": $gov_auth_token },
                    ($gov_priv_repo): { "access_token": $gov_auth_token },
                    "git@github.com:petar/gitrules.public.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" },
                    "git@github.com:petar/gitrules.private.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" }
               },
               "gov_public_url": $gov_pub_repo,
               "gov_public_branch": "main",
               "gov_private_url": $gov_priv_repo,
               "gov_private_branch": "main",
               "member_public_url": "git@github.com:petar/gitrules.public.git",
               "member_public_branch": "main",
               "member_private_url": "git@github.com:petar/gitrules.private.git",
               "member_private_branch": "main"
          }'
)
echo $CONFIG_JSON

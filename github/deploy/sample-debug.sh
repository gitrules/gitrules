#!/bin/sh

set -x -e

go install ../../gitrules
# ./sample-erase.sh
./sample-deploy.sh

gitrules --config sample-config.json user add --name petar --repo https://github.com/petar/gitrules-identity-public.git --branch main
gitrules --config sample-config.json account issue --to user:petar --asset plural -q 11000
gitrules --config sample-config.json account issue --to pmp+matching --asset plural -q 1000

# gitrules --config sample-config.json ballot vote --name pmp/motion/priority_poll/14 --choices rank --strengths 10.0
# gitrules --config sample-config.json ballot vote --name pmp/motion/approval_poll/13 --choices rank --strengths 20.0
# gitrules --config sample-config.json ballot show --name pmp/motion/approval_poll/13
# gitrules --config sample-config.json motion show --name 10

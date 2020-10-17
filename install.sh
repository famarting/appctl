#!/bin/sh
set -e

curl https://raw.githubusercontent.com/famartinrh/appctl/master/appctl --output /tmp/appctl-tmp
sudo chmod +x /tmp/appctl-tmp
mv /tmp/appctl-tmp ~/bin/appctl
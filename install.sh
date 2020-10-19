#!/bin/sh
set -e

curl https://raw.githubusercontent.com/famartinrh/appctl/master/appctl --output ~/bin/appctl
chmod +x ~/bin/appctl
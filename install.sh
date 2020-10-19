#!/bin/sh
set -e

curl https://raw.githubusercontent.com/famartinrh/appctl/master/appctl --output $HOME/bin/appctl
chmod +x $HOME/bin/appctl
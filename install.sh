#!/bin/sh
set -e

INSTALL_DIR=${INSTALL_DIR:-$HOME/bin}

curl https://raw.githubusercontent.com/famartinrh/appctl/master/appctl --output $INSTALL_DIR/appctl
chmod +x $INSTALL_DIR/appctl

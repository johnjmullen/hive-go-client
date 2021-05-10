#!/bin/bash
DEB_VERSION=$(echo 8.3.0 | sed -e 's|\.\([[:digit:]]\+\)-dev\.|-\1build|g' | sed -e 's|\.\([[:digit:]]\+\)$|-\1|g')
USERNAME=$(git config user.name)
TIMESTAMP=$(date +'%a, %d %b %Y %H:%M:%S %z')
DEBEMAIL=package@hiveio.com
DEBFULLNAME=HiveIO
dch --distribution unstable --package "hive-linux-agent" --newversion $DEB_VERSION "$DEB_VERSION release"

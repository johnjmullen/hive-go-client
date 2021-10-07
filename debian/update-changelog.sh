#!/bin/bash
DEB_VERSION=$(echo 8.5.0 | sed -e 's|\.\([[:digit:]]\+\)-dev\.|-\1build|g' | sed -e 's|\.\([[:digit:]]\+\)$|-\1|g')
DEB_VERSION="${DEB_VERSION}-$(date +%y%m%d%H)"
USERNAME=$(git config user.name)
DEBEMAIL=package@hiveio.com
DEBFULLNAME=HiveIO
dch --distribution unstable --package "hive-go-client" --newversion $DEB_VERSION "$DEB_VERSION release"

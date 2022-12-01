#!/bin/bash
DEB_VERSION="8.5.0-$(date +%y%m%d%H%M)"
USERNAME=$(git config user.name)
DEBEMAIL=package@hiveio.com
DEBFULLNAME=HiveIO
dch --distribution unstable --package "hive-go-client" --newversion $DEB_VERSION "$DEB_VERSION release"

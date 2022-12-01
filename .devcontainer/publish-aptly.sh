#!/bin/bash -ex
WORK_DIR=$1
REPO=$2
if [ x$(curl -s -o /dev/null -w "%{http_code}" "http://10.19.101.12:8082/api/repos/${REPO}") == x404 ];
then
  curl -X POST -H 'Content-Type: application/json' --data "{\"Name\": \"${REPO}\", \"DefaultDistribution\": \"bionic\", \"DefaultComponent\": \"main\"}" http://10.19.101.12:8082/api/repos || exit 1
fi
for PKG in $(ls /tmp/*.deb)
do
  NAME=$(dpkg-deb -f $PKG Package) || exit 1
  echo "found $NAME" || exit 1
  curl -X POST -F file=@$PKG http://10.19.101.12:8082/api/files/${NAME} && \
  curl -X POST http://10.19.101.12:8082/api/repos/${REPO}/file/${NAME} || exit 1
done
sleep 1
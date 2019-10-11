#!/bin/bash

VERSION=`git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $1;exit}'`

# clone
git clone --depth=50 --branch=master https://github.com/sacloud/docker-machine-sakuracloud-docker.git docker-machine-sakuracloud-docker
cd docker-machine-sakuracloud-docker
git fetch origin

# check version
CURRENT_VERSION=`git tag -l --sort=-v:refname | perl -ne 'if(/^([0-9\.]+)$/){print $1;exit}'`
if [ "$CURRENT_VERSION" = "$VERSION" ] ; then
    echo "docker-machine-sakuracloud-docker v$VERSION is already released."
    exit 0
fi

cat << EOL > Dockerfile
FROM alpine:3.10
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN set -x && apk add --no-cache curl ca-certificates zip && \
    curl -L https://github.com/docker/machine/releases/download/v${DOCKER_MACHINE_VERSION}/docker-machine-Linux-x86_64 >/usr/local/bin/docker-machine && \
    chmod +x /usr/local/bin/docker-machine && \
    curl -LO https://github.com/sacloud/docker-machine-sakuracloud/releases/download/v${VERSION}/docker-machine-driver-sakuracloud_linux-amd64.zip && \
    unzip docker-machine-driver-sakuracloud_linux-amd64.zip -d /usr/local/bin/ && \
    chmod +x /usr/local/bin/docker-machine-driver-sakuracloud

VOLUME ["/workdir"]
WORKDIR /workdir

ENTRYPOINT ["/usr/local/bin/docker-machine"]
CMD ["--help"]

EOL

git config --global push.default matching
git config user.email 'sacloud.users@gmail.com'
git config user.name 'sacloud-bot'
git commit -am "v${VERSION}"
git tag "${VERSION}"

echo "Push ${VERSION} to github.com/sacloud/docker-machine-sakuracloud-docker.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/docker-machine-sakuracloud-docker.git" >& /dev/null

echo "Cleanup tag ${VERSION} on github.com/sacloud/docker-machine-sakuracloud-docker.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/docker-machine-sakuracloud-docker.git" :${VERSION} >& /dev/null

echo "Tagging ${VERSION} on github.com/sacloud/docker-machine-sakuracloud-docker.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/docker-machine-sakuracloud-docker.git" ${VERSION} >& /dev/null
exit 0

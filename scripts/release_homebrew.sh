#!/bin/bash

VERSION=`git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $1;exit}'`
SHA256_SRC_DARWIN=`openssl dgst -sha256 bin/docker-machine-driver-sakuracloud_darwin-amd64.zip | awk '{print $2}'`
SHA256_SRC_LINUX=`openssl dgst -sha256 bin/docker-machine-driver-sakuracloud_linux-amd64.zip | awk '{print $2}'`
# clone
git clone --depth=50 --branch=master https://github.com/sacloud/homebrew-docker-machinbe-sakuracloud.git homebrew-docker-machine-sakuracloud
cd homebrew-docker-machine-sakuracloud

# check version
CURRENT_VERSION=`git log --oneline | perl -ne 'if(/^.+ v([0-9\.]+)/){print $1;exit}'`
if [ "$CURRENT_VERSION" = "$VERSION" ] ; then
    echo "homebrew-docker-machine-sakuracloud v$VERSION is already released."
    exit 0
fi

cat << EOL > docker-machine-sakuracloud.rb
class DockerMachineSakuracloud < Formula

  _version = "${VERSION}"
  sha256_src_darwin = "${SHA256_SRC_DARWIN}"
  sha256_src_linux = "${SHA256_SRC_LINUX}"

  desc "Docker-Machine driver for SakuraCloud"
  homepage "https://github.com/sacloud/docker-machine-sakuracloud"
  head "https://github.com/sacloud/docker-machine-sakuracloud.git"
  version _version

  if OS.mac?
    url "https://github.com/sacloud/docker-machine-sakuracloud/releases/download/v#{_version}/docker-machine-driver-sakuracloud_darwin-amd64.zip"
    sha256 sha256_src_darwin
  else
    url "https://github.com/sacloud/docker-machine-sakuracloud/releases/download/v#{_version}/docker-machine-driver-sakuracloud_linux-amd64.zip"
    sha256 sha256_src_linux
  end

  depends_on "docker-machine" => :run

  def install
    bin.install "docker-machine-driver-sakuracloud"
  end

  test do
    assert_match "sakuracloud-access-token", shell_output("docker-machine create -d sakuracloud -h")
  end
end
EOL

git config --global push.default matching
git config user.email 'sacloud.users@gmail.com'
git config user.name 'sacloud-bot'
git commit -am "v${VERSION}"

echo "Push ${VERSION} to github.com/sacloud/homebrew-terraform-provider-sakuracloud.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-docker-machine-sakuracloud.git" >& /dev/null

echo "Cleanup tag v${VERSION} on github.com/sacloud/homebrew-docker-machine-sakuracloud.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-docker-machine-sakuracloud.git" :v${VERSION} >& /dev/null

echo "Tagging v${VERSION} on github.com/sacloud/homebrew-docker-machine-sakuracloud.git"
git tag v${VERSION} >& /dev/null
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-docker-machine-sakuracloud.git" v${VERSION} >& /dev/null
exit 0

# Docker Machine SAKURA CLOUD Driver

[Docker Machine](https://docs.docker.com/machine/)で[さくらのクラウド](http://cloud.sakura.ad.jp)上の仮想マシンを利用できるようにするプラグインです。

([English version](README.en.md))

[![Build Status](https://travis-ci.org/sacloud/docker-machine-sakuracloud.svg?branch=master)](https://travis-ci.org/sacloud/docker-machine-sakuracloud)

## Quick Start: Dockerコンテナでの実行

docker-machineとdocker-machine-sakuracloud同梱したDockerイメージを用意しています。

[sacloud/docker-machine](https://hub.docker.com/r/sacloud/docker-machine/)

お手元のマシンにDockerが無い場合、以降の[インストール](#インストール)を参照ください。

Dockerでの実行方法

```bash:書式
docker run [dockerコマンドオプション] sacloud/docker-machine [docker-machineオプション] <作成するマシン名>
```

マシンの作成を行う場合、以下のように実行します。

```bash:コマンド例
docker run -it --rm -e MACHINE_STORAGE_PATH=$HOME/.docker/machine \
                    -e SAKURACLOUD_ACCESS_TOKEN=<トークン> \
                    -e SAKURACLOUD_ACCESS_TOKEN_SECRET=<シークレット> \
                    -v $HOME/.docker:$HOME/.docker \
                    sacloud/docker-machine create -d sakuracloud sakura-dev
```

## インストール

### macOS(HomeBrew) または Linux(LinuxBrew)

    $ brew tap sacloud/docker-machine-sakuracloud
    $ brew install docker-machine-sakuracloud

### マニュアルインストール

あらかじめ`docker-machine`をインストールしておいてください。

`docker-machine-driver-sakuracloud`バイナリをダウンロードし、パス(`$PATH`)を通してください。
(Windowsの場合はdocker-machine.exeと同じフォルダに配置すればよいです)
配置後にchmod +xしておいてください。

`docker-machine-driver-sakuracloud`の最新のバイナリはこちらの["Releases"](https://github.com/sacloud/docker-machine-sakuracloud/releases/latest)ページからダウンロードできます。

## 使い方

Docker Machine 公式ドキュメントは[こちら](https://docs.docker.com/machine/)。
Windowsの場合は以下のページも参照してください。
 - [Instration on Windows](http://docs.docker.com/engine/installation/windows/)。


さくらのクラウドのコントロールパネルからAPIキーを発行しておいてください。

さくらのクラウド上にdocker-machineで利用する仮想マシンを作成するには以下のコマンドを実施してください。

```
$ docker-machine create --driver=sakuracloud \
    --sakuracloud-access-token=[アクセストークン] \
    --sakuracloud-access-token-secret=[アクセストークンシークレット] \
    sakura-dev
```

オプション:

 - `--sakuracloud-access-token`: **必須** アクセストークン
 - `--sakuracloud-access-token-secret`: **必須** アクセストークンシークレット
 - `--sakuracloud-zone`: 対象ゾーン[`is1a` / `is1b` / `tk1a`]
 - `--sakuracloud-os-type`: OS[`rancheros` / `centos` / `ubuntu`]
 - `--sakuracloud-core`: CPUコア数
 - `--sakuracloud-memory`: メモリサイズ(GB単位)
 - `--sakuracloud-disk-connection`: ディスクインターフェース (`virtio` or `ide`)
 - `--sakuracloud-disk-plan`: ディスクプラン (`ssd` / `hdd`)
 - `--sakuracloud-disk-size`: ディスクサイズ(GB単位)
 - `--sakuracloud-interface-driver`: NICドライバ(`virtio` or `e1000`)
 - `--sakuracloud-password`: 管理ユーザーのパスワード(未指定の場合ランダムな文字列を利用)
 - `--sakuracloud-enable-password-auth` : SSHでのパスワード認証の有効化(デフォルトは公開鍵認証のみが有効)
 - `--sakuracloud-packet-filter`: パケットフィルタのID
 - `--sakuracloud-engine-port` : Docker Engineのポート番号
 - `--sakuracloud-ssh-key` : SSH秘密鍵へのパス(省略した場合は新たなキーペアが生成されます)

`--sakuracloud-disk-size`はさくらのクラウドでサポートされるサイズのみ指定可能です。
サポートされるサイズについては[サービス仕様・料金](http://cloud.sakura.ad.jp/specification.php)ページを参照してください。
また、`--sakuracloud-disk-plan`の選択によってサポートされるサイズが変わるため注意してください。

`--sakuracloud-zone`では利用したいリージョンに応じて以下の値を指定してください。
SandboxリージョンについてはSSHにてログインができないため利用できません。

 - 石狩第1ゾーン : `is1a`
 - 石狩第2ゾーン : `is1b`
 - 東京第1ゾーン : `tk1a`


各オプションは環境変数で指定することも可能です。


環境変数名とデフォルト値:

| CLI option                           | Environment variable              | Default                  |
|--------------------------------------|-----------------------------------|--------------------------|
| `--sakuracloud-access-token`         | `SAKURACLOUD_ACCESS_TOKEN`        | -                        |
| `--sakuracloud-access-token-secret`  | `SAKURACLOUD_ACCESS_TOKEN_SECRET` | -                        |
| `--sakuracloud-zone`                 | `SAKURACLOUD_ZONE`                | `is1b`                   |
| `--sakuracloud-os-type`              | `SAKURACLOUD_OS_TYPE`             | `rancheros`              |
| `--sakuracloud-core`                 | `SAKURACLOUD_CORE`                | `1`                      |
| `--sakuracloud-memory`               | `SAKURACLOUD_MEMORY`              | `1`                      |
| `--sakuracloud-disk-connection`      | `SAKURACLOUD_DISK_CONNECTION`     | `virtio`                 |
| `--sakuracloud-disk-plan`            | `SAKURACLOUD_DISK_PLAN`           | `ssd`                    |
| `--sakuracloud-disk-size`            | `SAKURACLOUD_DISK_SIZE`           | `20`                     |
| `--sakuracloud-interface-driver`     | `SAKURACLOUD_INTERFACE_DRIVER`    | `virtio`                 |
| `--sakuracloud-password`             | `SAKURACLOUD_PASSWORD`            | -                        |
| `--sakuracloud-enable-password-auth` | `SAKURACLOUD_ENABLE_PASSWORD_AUTH`| false                    |
| `--sakuracloud-packet-filter`        | `SAKURACLOUD_PACKET_FILTER`       | -                        |
| `--sakuracloud-engine-port`          | `SAKURACLOUD_ENGINE_PORT`         | `2376`                   |
| `--sakuracloud-ssh-key`              | `SAKURACLOUD_SSH_KEY`             | -                        |

## Author

* Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))

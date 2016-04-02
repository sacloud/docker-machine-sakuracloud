# Docker Machine SAKURA CLOUD Driver

[Docker Machine](https://docs.docker.com/machine/)で[さくらのクラウド](http://cloud.sakura.ad.jp)上の仮想マシンを利用できるようにするプラグインです。

docker-machine 0.5にて[マージ](https://github.com/docker/machine/pull/1902)された、外部プラグイン機構に対応しています。

[![Build Status](https://travis-ci.org/yamamoto-febc/docker-machine-sakuracloud.svg?branch=master)](https://travis-ci.org/yamamoto-febc/docker-machine-sakuracloud)

## 動作環境
* [Docker Machine](https://docs.docker.com/machine/) 0.5.1+ (is bundled to
  [Docker Toolbox](https://www.docker.com/docker-toolbox) 1.9.1+)

## 動作確認済み環境
* OSX 10.9+  : amd64
* Windows 10 : amd64

## インストール

#### Windowsの場合:

[こちら](https://github.com/yamamoto-febc/docker-machine-sakuracloud/releases/download/v0.0.8/DockerMachineSakuracloudSetup.exe)からインストーラーをダウンロードして実行してください。

#### OSX(Mac)の場合:

Homebrewでインストールします。以下のコマンドを実行してください。

```console
$ brew tap yamamoto-febc/docker-machine-sakuracloud
$ brew install docker-machine-sakuracloud
```
#### 手動でインストールする場合:

`docker-machine-driver-sakuracloud`バイナリをダウンロードし、パス(`$PATH`)を通してください。
(Windowsの場合はdocker-machine.exeと同じフォルダに配置すればよいです)
配置後にchmod +xしておいてください。

```console
$ chmod +x /usr/local/bin/docker-machine-driver-sakuracloud
```

`docker-machine-driver-sakuracloud`の最新のバイナリはこちらの["Releases"](https://github.com/yamamoto-febc/docker-machine-sakuracloud/releases/latest)ページからダウンロードできます。

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
 - `--sakuracloud-core`: CPUコア数
 - `--sakuracloud-memory-size`: メモリサイズ(GB単位).
 - `--sakuracloud-connected-switch`: 接続するスイッチ or ルーターのID(eth1)
 - `--sakuracloud-private-ip-only`: 公開セグメントのNICを無効にしeth1のみを使う
 - `--sakuracloud-private-ip`: eth1のIPアドレス
 - `--sakuracloud-private-subnet-mask`: eth1のサブネットマスク
 - `--sakuracloud-gateway`: デフォルトゲートウェイ(eth1を使う場合は必須)
 - `--sakuracloud-disk-connection`: ディスクインターフェース (`virtio` or `ide`)
 - `--sakuracloud-disk-name`: さくらのクラウド上に作成するディスクの名前
 - `--sakuracloud-disk-plan`: ディスクプラン (HDD:`2` or SSD:`4`)
 - `--sakuracloud-disk-size`: ディスクサイズ(MB単位)
 - `--sakuracloud-region`: リージョン名[is1a / is1b / tk1a]
 - `--sakuracloud-group`: @group 特殊タグ
 - `--sakuracloud-auto-reboot`: @auto-reboot 特殊タグ
 - `--sakuracloud-ignore-virtio-net`: 順仮想化ドライバの無効化(@virtio-net-pci 特殊タグの無効化)
 - `--sakuracloud-packet-filter`: パケットフィルタのID/名称 (eth0:共有セグメント用)
 - `--sakuracloud-private-packet-filter`: パケットフィルタのID/名称 (eth1:プライベートセグメント用)
 - `--sakuracloud-enable-password-auth` : SSHでのパスワード認証の有効化(デフォルトは公開鍵認証のみが有効)


`--sakuracloud-disk-size`はさくらのクラウドでサポートされるサイズのみ指定可能です。
サポートされるサイズについては[サービス仕様・料金](http://cloud.sakura.ad.jp/specification.php)ページを参照してください。(1GB = 1024MB)
また、`--sakuracloud-disk-plan`の選択によってサポートされるサイズが変わるため注意してください。


`--sakuracloud-region`では利用したいリージョンに応じて以下の値を指定してください。
Sandboxリージョンについては外部からログインができないため利用できません。

 - 石狩第1ゾーン : `is1a`
 - 石狩第2ゾーン : `is1b`
 - 東京第1ゾーン : `tk1a`


各オプションは環境変数で指定することも可能です。


環境変数名とデフォルト値:

| CLI option                          | Environment variable              | Default                  |
|-------------------------------------|-----------------------------------|--------------------------|
| `--sakuracloud-access-token`        | `SAKURACLOUD_ACCESS_TOKEN`        | -                        |
| `--sakuracloud-access-token-secret` | `SAKURACLOUD_ACCESS_TOKEN_SECRET` | -                        |
| `--sakuracloud-core`                | `SAKURACLOUD_CORE`                | `1`                   |
| `--sakuracloud-memory-size`         | `SAKURACLOUD_MEMORY_SIZE`         | `1`                   |
| `--sakuracloud-connected-switch`    | `SAKURACLOUD_CONNECTED_SWITCH`     | -                 |
| `--sakuracloud-private-ip-only`       | `SAKURACLOUD_PRIVATE_IP_ONLY`     | -                 |
| `--sakuracloud-private-ip`       | `SAKURACLOUD_PRIVATE_IP`     | -                 |
| `--sakuracloud-private-ip-subnet-mask`     | `SAKURACLOUD_PRIVATE_IP_SUBNET_MASK`     | `255.255.255.0`                 |
| `--sakuracloud-gateway`     | `SAKURACLOUD_GATEWAY`     | -                 |
| `--sakuracloud-disk-connection`     | `SAKURACLOUD_DISK_CONNECTION`     | `virtio`                 |
| `--sakuracloud-disk-name`           | `SAKURACLOUD_DISK_NAME`           | `disk001`                |
| `--sakuracloud-disk-plan`           | `SAKURACLOUD_DISK_PLAN`           | `4`                      |
| `--sakuracloud-disk-size`           | `SAKURACLOUD_DISK_SIZE`           | `20480`                  |
| `--sakuracloud-region`              | `SAKURACLOUD_REGION`              | `is1a`                   |
| `--sakuracloud-group`               | `SAKURACLOUD_GROUP`              | -                   |
| `--sakuracloud-auto-reboot`         | `SAKURACLOUD_AUTO_REBOOT`        | -                   |
| `--sakuracloud-ignore-virtio-net`   | `SAKURACLOUD_IGNORE_VIRTIO_NET`  | -                   |
| `--sakuracloud-upgrade-kernel`   | `SAKURACLOUD_UPGRADE_KERNEL`  | -                   |
| `--sakuracloud-enable-password-auth`   | `SAKURACLOUD_ENABLE_PASSWORD_AUTH`  | false                 |


## Author

* Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))

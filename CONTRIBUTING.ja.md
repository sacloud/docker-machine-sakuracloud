# コントリビューションガイド:docker-machine-sakuracloud

([English version](CONTRIBUTING.md))

docker-machine-sakuracloudは皆様のご協力を **大歓迎** します。
Issueの投稿、Pull Requestsなどいつでもお待ちしています。
以下の簡単なルールだけ目を通しておいてください。

## Issueの投稿、Pull Requestの前に

 - 類似のIssueがないか検索を行ってください。
 - 類似のPull Requestの検索を行ってください。
 - Issue,Pull Requestとも本文はできれば日本語 or 英語で記載してください。

その他の細かい点はコメントでやりとりすれば良いと思います。

## docker-machine-sakuracloud開発環境について

ここではdocker-machine-sakuracloud開発の始め方を説明します。


# ビルド

docker-machine-sakuracloudのビルドには以下が必要です。

1. 実行中のDockerインスタンス(デーモン) or Golang 1.5 開発環境
2. `bash` シェル (MinGWもOK)
3. [Make](https://www.gnu.org/software/make/)

## Dockerコンテナでのビルド

Dockerコンテナで`docker-machine-sakuracloud`バイナリをビルドする場合、次のコマンドを実行します。

    $ export USE_CONTAINER=true
    $ make build

*v0.0.3以降、docker-machine-sakuracloudで作成したさくらのクラウド上のサーバで
dockerコンテナを用いてビルド可能となりました。*

## ローカルのGo1.5開発環境でビルド

新たにdocker-machine-driver-sakuracloudディレクトリを作成し、
GOPATHの設定、ソースの取得、依存ライブラリのインストールを行います。

```
    $ mkdir docker-machine-sakuracloud
    $ cd docker-machine-sakuracloud
    $ export GOPATH="$PWD"
    $ go get github.com/yamamoto-febc/docker-machine-sakuracloud
    $ cd src/github.com/yamamoto-febc/docker-machine-sakuracloud
    $ go get github.com/tools/godep
    $ $GOPATH/bin/godep get
```

環境設定ができたら以下コマンドを実行します。

    $ make build

## ビルド成果物

ビルドが成功すると`bin/docker-machine-driver-sakuracloud`バイナリが作成されています。

以下コマンドでビルド成果物を削除します。

    $ make clean

### ビルドターゲットなど

サポートされている全てのOS/CPUアーキテクチャ用のバイナリのビルドを行う場合、
次のコマンドを実行します(成果物はOS/CPU別に/binのサブディレクトリ配下に生成されます)。

    make cross

OS/CPUアーキテクチャを指定する場合、以下のようにします。

    TARGET_OS=linux TARGET_ARCH="amd64 arm" make cross

さらに、以下の環境変数を設定することでビルドオプションを設定可能です。

    DEBUG=true # enable debug build
    STATIC=true # build static (note: when cross-compiling, the build is always static)
    VERBOSE=true # verbose output
    PREFIX=folder # put binaries in another folder (not the default `./bin`)

### Go依存ライブラリの保存/復元

以下コマンドにてgodepを用いたGo依存ライブラリの管理が可能です。

    make dep-save
    make dep-restore

*makeを通じてgodepを実行すると`GO15VENDOREXPERIMENT`環境変数が設定されます。
これにより`vendor`ディレクトリ配下に依存ライブラリがインストールされます*
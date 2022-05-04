# atc

Atcoderのコマンドラインツール。

## Feature

- コンテストごとのディレクトリを作成できる
- テストや提出をCLI上から行うことができる
- 特定の言語に依存していない

## Install

```shell
go install github.com/DuGlaser/atc@latest
```

## Quick Start

### 1. Atcoderにログインする

以下のコマンドを実行し、Atcoderにログインします。

```shell
atc login
```

このコマンドにより、cacheディレクトリに`atc/session.txt`が作成され、Cookieが保存されます。

### 2. atcの設定をする

以下のコマンドを実行し、atcの設定を行います。

```shell
atc config
```

このコマンドを実行することで、configディレクトリに`.atc.toml`が作成され、設定が保存されます。

### 3. コンテストを作成する

以下のコマンドを実行し、コンテスト用のディレクトリを作成することができます。

```sholl
atc new abc001
```

このコマンドを実行することで、`abc001`用のディレクトリを作成すｒことができます。
`abc001`となっているところは好きなコンテストに変えることができます。

### 4. 問題のテストコードを実行する

まずは先程作成した`abc001`用のディレクトリに移動します。

```shell
cd abc001/
```

次に以下のコマンドを実行し、A問題のテストケースを実行します。

```shell
atc test a
```

すると以下のようにすべてのテストケースが通らなかったことがわかります。

```shell
sample test case 1 ... failed
sample test case 2 ... failed
sample test case 3 ... failed

=== sample test case 1 ===

input:
 1 | 15
 2 | 10

expected:
 1 | 5

your output:
 1 |

=== sample test case 2 ===

input:
 1 | 0
 2 | 0

expected:
 1 | 0

your output:
 1 |

=== sample test case 3 ===

input:
 1 | 5
 2 | 20

expected:
 1 | -15

your output:
 1 |
```

正しいコードを`a`ディレクトリ以下のファイルに記述し、再度実行するとすべてのテストケースが成功したことがわかります。

```shell
ample test case 1 ... success
sample test case 2 ... success
sample test case 3 ... success
```

### 5. 解答を提出する

次に以下のコマンドを実行し、A問題の解答を提出します。

```shell
atc submit a
```

submitを行う際は先程テストを行ったのと同様に、テストケースを自動的に実行し、全て成功したときのみAtcoderに提出します。


## Usage

```
Atcoder command line tool

Usage:
  atc [command]

Available Commands:
  browse      open problem in browse
  completion  Generate the autocompletion script for the specified shell
  config      Create and edit atc config
  help        Help about any command
  info        Output the information of logged in users
  login       Login to atcoder
  logout      Logout to atcoder
  new         Create contest project
  submit      Submit answer
  test        Test answer

Flags:
      --config string   config file (default is $HOME/.config/.atc.toml)
  -h, --help            help for atc
  -t, --toggle          Help message for toggle
  -v, --verbose         Make the operation more talkative

Use "atc [command] --help" for more information about a command.
```

### browse

コンテストの問題をブラウザで開くコマンドです。
このコマンドはコンテストディレクトリのみで使用することができます。

### config

atcの設定を生成することができます。

```toml
[config]
  runcmd = "" # テストを実行するときに実行するコマンド。
  buildcmd = "" # テストを実行するときに一度だけ実行されるコマンド。必ずruncmdよりも前に実行されます。
  filename = "" # 実際に解答を記述するファイル名。
  lang = "" # Atcoderに提出する言語のID。
  template = """
  """ # 提出するファイルのテンプレート。
```


#### runcmd

テストを実行するときに実行するコマンドです。
このコマンドには以下のような実行時に埋め込まれる値を指定することができます。

- `{{ dir }}`: 実行するファイルが存在するディレクトリのpathが埋め込まれます。
- `{{ file }}`: 実行するファイルのpathが埋め込まれます。


#### buildcmd

テストを実行するときに一度だけ実行されるコマンドです。必ずruncmdよりも前に実行されます。
buildcmdを使用しない場合は空文字列を指定します。
このコマンドには以下のような実行時に埋め込まれる値を指定することができます。

- `{{ dir }}`: 実行するファイルが存在するディレクトリのpathが埋め込まれます。
- `{{ file }}`: 実行するファイルのpathが埋め込まれます。

#### filename

実際に解答を記述するファイル名を指定します。

#### template

コンテストディレクトリを作成したときに解答ファイルに記述しておく内容を指定することができます。

### login

Atcoderにloginする際に使用します。
このコマンドにより、cacheディレクトリに`atc/session.txt`が作成され、Cookieが保存されます。

### logout

cacheディレクトリに`atc/session.txt`を削除するコマンドです。

### new

コンテスト用のディレクトリを作成するためのコマンドです。

### submit

コンテストの問題を提出するコマンドです。
このコマンドはコンテストディレクトリのみで使用することができます。

### test

コンテストの問題のテストケースを実行するコマンドです。
このコマンドはコンテストディレクトリのみで使用することができます。

## Example

### C++

```toml
[config]
  runcmd = "{{ .dir }}/main"
  buildcmd = "g++ -o {{ .dir }}/main {{ .file }}"
  filename = "main.cpp"
  lang = "4003"
  template = """
#include <bits/stdc++.h>
using namespace std;

int main() {

}
"""
```

### Go

```toml
[config]
  runcmd = "go run {{ .file }}"
  buildcmd = ""
  filename = "main.go"
  lang = "4026"
  template = """
package main

func main() {

}
"""
```

### Python

```toml
[config]
  runcmd = "python3 {{ .file }}"
  buildcmd = ""
  filename = "main.py"
  lang = "4006"
  template = """
"""
```

## LICENSE

[MIT](https://github.com/DuGlaser/atc/blob/master/LICENSE)

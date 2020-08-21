## `go generate`検証用リポジトリ

### 検証項目
1. 自作のコードジェネレーターの一番簡単な指定方法
1. 一つのファイルから伝搬的に複数のファイルを生成するときの挙動
   - 例) AがBを生成してBがCを生成する

### 結果
- まずは自作のコードジェネレーターを1ファイル受け取って静的解析をして別ファイルを作るツールにしなくてはならない
  - 今まで作ったものがimportしたデータからreflectを使って動的解析するパターンが多かったので静的解析を学ばなければ
- goの相対パス周りでかなり躓いた
  - template.ParseFileを使うのを諦めた
- 検証項目2は無理？
  - 例のような場合、AからBを生成したあとCを作るにはもう一度コマンドを実行する必要がある
    - 生成後にAからBに対して`go generate B`すればいける？
        - いけた
- template上にある`//go:generate`が反応してしまう・・・
  - `cmd/generate/tpl/mock_tpl`が生成されてしまった
  - `go generate ./...`ではなく`go generate ./pkg`等でcmdは除外する

```
 $ go generate -v -x ./...
cmd/generator/main.go
cmd/generator/tpl/template.go
mockgen -source=template.go -destination=mock_tpl/mock_template.go
goimports -w --local github.com/iyu/go-generator-test mock_tpl/mock_template.go
cmd/generator/tpl/mock_tpl/mock_template.go
pkg/domain/entity/user/entity.go
go run ../../../../cmd/generator entity.go ../../repository/user/repository.go
go generate ../../repository/user/repository.go
```

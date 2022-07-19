# Protocol Buffersとは
Googleによって2008年にオープンソース化されたスキーマ言語

【スキーマ言語】
何かしらの処理をさせるのではなく、要素や属性などの構造を定義するための言語

## スキーマ言語がなぜ重要か
マイクロサービス化が主流となっており、PCだけでなくIOSやAndroidのクライアント対応も必須となっている。
→ 事前にスキーマ言語で宣言的にI/Fを定義しておく

## Protocol Buffersの特徴
- gRPCのデータフォーマットとして使用されている

- プログラミング言語からは独立しており、様々な言語に変換可能
`Java, Go, Python, C++, Kotolin, etc...`


- バイナリ形式にシリアライズするので、サイズが小さく拘束な通信が可能

- 型安全にデータのやり取りが可能

- JSONに変換することも可能


## JSONとの比較
![Jsonとの比較](images/jsonとの比較.png)


## Protocol Buffersを使用した開発の進め方
1. スキーマの定義
1. 開発言語のオブジェクトを自動生成
1. バイナリ形式へシリアライズ
![ProtocolBuffersを使用した開発の進め方](/images/ProtocolBuffersを使用した開発の進め方.png)

# messageとは
- 複数のフィールドを持つことができる定義
- - それぞのフィールドはスカラ型もしくはコンポジット型
- 各言語のコードとしてコンパイルした場合、構造体やクラスとして変換される
- 1つのprotoファイルに複数のmessage型を定義することも可能

![message例](images/message例.png)
- メッセージ名;
- フィールド項目：型 - 名前 - タグ番号;

`※行末にはセミコロンが必要`

## スカラー型
[https://developers.google.com/protocol-buffers/docs/proto3#scalar](https://developers.google.com/protocol-buffers/docs/proto3#scalar)

## タグ
- Protocol Buffersではフィールドはフィールド名ではなく、タグ番号によって識別される
- タグ番号の重複は許されず、一意である必要がある
- タグの最小値は1, 最大値は2^29-1(値は536,870,911)
- 19000 ~ 19999はProtocol Buffersの予約番号のため使用不可

- 1~15番までは1byteで表すことができるため、よく使うフィールドは1~15番を割り当てる
(それを意識したメッセージにすることでパフォーマンス向上が望める)

- タグは連番にする必要はないので、あまり使わないフィールドはあえて16番以降を割り当てることも可能

- タグ番号を予約するなど、安全にProtocol Buffersを使用する方法も用意されている

## 列挙型
タグ番号は0から始まる、必ず0から始める必要がある
慣例的に0番はUNKNOWNにすることが多い

## デフォルト値
- 定義したmessageでデータをやり取りする際に、定義したフィールドがセット去れていない場合、そのフィールドのデフォルト値が設定される
- デフォルト値は型によって決められている

- - string: 空の文字列
- - bytes: 空のbyte
- - bool: false
- - 整数型・浮動小数点型: 0
- - 列挙型: タグ番号0の値
- - repeated: 空のリスト

## protocコマンド
- `-IPATH, --proto_path=PATH`
![IPATH1](images/IPATH1.png)
![IPATH2](images/IPATH2.png)

- `各言語に変換するためのオプション`
- - オプションによってどの言語に変換するかを決定する
- - `※Go言語のオプションはプラグインで追加する必要がある`
```
--cpp_out=OUT_DIR
--csharp_out=OUT_DIR
--java_out=OUT_DIR
--js_out=OUT_DIR
--kotlin_out=OUT_DIR
--objec_out=OUT_DIR
--php_out=OUT_DIR
--python_out=OUT_DIR
--ruby_out=OUT_DIR
--go_out=OUT_DIR ※プラグインが必要
```

- `コンパイルするファイルの指定`
```
protoc -I. --go_out=. proto/employee.proto proto/date.proto
```

```
protoc -I. --go_out=. proto/*.proto
```



# protoファイル生成コマンド(今回実行コマンド)
```
protoc -I. --go_out=. proto/*.proto

# pb/date.pb.goのimportエラーの為
go mod tidy


```






